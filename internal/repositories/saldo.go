package repositories

import (
	"context"
	"crypto/sha1"
	"database/sql"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Client represents Otomax API client
type Client struct {
	httpClient *http.Client
}
type ParsingSupplierBalance struct {
	MemberId string
	Pin      string
	Password string
}

const (
	DefaultTimeout = 30 * time.Second
	UserAgent      = "Go-Otomax-Client/1.0"
	AcceptHeader   = "text/plain,text/html,application/json,*/*"
)

func (repo *ProductRepository) GetBalanceSupplier(c context.Context, kode int) (string, error) {
	var parameter_parsing, tujuan string
	query := `
		select parameter_parsing,tujuan from modul 
		where kode = @kode_modul
	`
	err := repo.db.QueryRowContext(c, query, sql.Named("kode_modul", kode)).Scan(&parameter_parsing, &tujuan)

	if err != nil {
		return "", fmt.Errorf("failed to query balance: %w", err)
	}

	client := &Client{
		httpClient: &http.Client{
			Timeout: DefaultTimeout,
		},
	}

	parsed := RegexParsingSupplierBalance(parameter_parsing)
	balance, err := client.CheckBalance(c, tujuan, parsed.MemberId, parsed.Pin, parsed.Password)
	if err != nil {
		return "", fmt.Errorf("failed to check balance: %w", err)
	}

	return balance, nil
}

// makeRequest makes HTTP request with proper error handling
func (c *Client) makeRequest(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("Accept", AcceptHeader)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	return resp, nil
}

func RegexParsingSupplierBalance(parsing string) ParsingSupplierBalance {

	parts := strings.Split(parsing, "|")

	var result ParsingSupplierBalance
	for _, v := range parts {
		parts := strings.Split(v, "=")
		if len(parts) == 2 {
			key := parts[0]
			value := parts[1]
			switch key {
			case "memberId":
				result.MemberId = value
			case "pin":
				result.Pin = value
			case "password":
				result.Password = value
			}
		}
	}
	return result
}

// CheckBalance checks account balance
func (c *Client) CheckBalance(ctx context.Context, baseUrl, memberID, pin, password string) (string, error) {
	url := c.buildBalanceURL(baseUrl, memberID, pin, password)

	fmt.Println("Balance URL:", url) // Debug print
	resp, err := c.makeRequest(ctx, url)
	if err != nil {
		return "", err
	}

	bodyStr, err := c.readResponseBody(resp)
	if err != nil {
		return "", err
	}

	return bodyStr, nil
}
func (c *Client) readResponseBody(resp *http.Response) (string, error) {
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	return string(body), nil
}

func (c *Client) buildBalanceURL(baseUrl, memberID, pin, password string) string {
	sign := SignCheckBalanceOtomax(memberID, pin, password)
	baseUrl = strings.TrimRight(baseUrl, "/")
	return fmt.Sprintf("%s/balance?memberID=%s&sign=%s", baseUrl, memberID, sign)
}

func SignCheckBalanceOtomax(memberID, pin, password string) string {
	payload := fmt.Sprintf("OtomaX|CheckBalance|%s|%s|%s", memberID, pin, password)

	// Create SHA1 hash
	hasher := sha1.New()
	hasher.Write([]byte(payload))
	hashBytes := hasher.Sum(nil)

	// Convert to base64
	base64Hash := base64.StdEncoding.EncodeToString(hashBytes)
	urlSafe := strings.ReplaceAll(base64Hash, "+", "-")
	urlSafe = strings.ReplaceAll(urlSafe, "/", "_")
	urlSafe = strings.ReplaceAll(urlSafe, "=", "")

	return urlSafe
}
