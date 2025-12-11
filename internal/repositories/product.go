package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/universeh2h/report/internal/model"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{
		db: db,
	}
}

func (repo *ProductRepository) TotalTransaksi(c context.Context, req model.PaginationParams) (*model.AnalyticsResult, error) {
	// Build WHERE conditions
	whereConditions := ""
	args := []interface{}{}

	if req.Search != nil && *req.Search != "" {
		whereConditions += " AND kode_reseller LIKE @search"
		args = append(args, sql.Named("search", "%"+*req.Search+"%"))
	}

	if req.StartDate != "" {
		whereConditions += " AND CAST(tgl_status AS DATE) >= @startdate"
		args = append(args, sql.Named("startdate", req.StartDate))
	}

	if req.EndDate != "" {
		whereConditions += " AND CAST(tgl_status AS DATE) <= @enddate"
		args = append(args, sql.Named("enddate", req.EndDate))
	}

	query := `
		SELECT 
    'TOTAL' as type,
    '' as kode_reseller,
    '' as kode_produk,
    COUNT(*) as jumlah_transaksi,
    COALESCE(SUM(harga - harga_beli), 0) as total_laba,
    '' as kode_operator,
    '' as nama,
    0 as sort_order
FROM transaksi 
WHERE status = 20` + whereConditions + `

UNION ALL

SELECT 
    'RESELLER' as type,
    t.kode_reseller,
    '' as kode_produk,
    COUNT(*) as jumlah_transaksi,
    COALESCE(SUM(t.harga - t.harga_beli), 0) as total_laba,
    '' as kode_operator,
    r.nama,
    1 as sort_order
FROM transaksi t
LEFT JOIN reseller r ON t.kode_reseller = r.kode
WHERE t.status = 20` + whereConditions + `
GROUP BY t.kode_reseller, r.nama

UNION ALL

SELECT 
    'PRODUK' as type,
    t.kode_reseller,
    t.kode_produk,
    COUNT(*) as jumlah_transaksi,
    COALESCE(SUM(t.harga - t.harga_beli), 0) as total_laba,
    p.kode_operator,
    '' as nama,
    2 as sort_order
FROM transaksi t
LEFT JOIN produk p ON t.kode_produk = p.kode
WHERE t.status = 20` + whereConditions + `
GROUP BY t.kode_reseller, t.kode_produk, p.kode_operator

ORDER BY sort_order, kode_reseller ASC, kode_produk ASC
	`

	rows, err := repo.db.QueryContext(c, query, args...)
	if err != nil {
		log.Printf("failed to execute nested analytics query: %v", err)
		return nil, fmt.Errorf("failed to execute nested analytics query: %w", err)
	}
	defer rows.Close()

	var result model.AnalyticsResult
	resellerMap := make(map[string]*model.ResellerAnalytics)
	var resellerOrder []string

	for rows.Next() {
		var tipe, kodeReseller, kodeProduk string
		var jumlahTrx int
		var totalLaba float64
		var sortOrder int
		var NamaReseller, kodeOperator string

		err := rows.Scan(&tipe, &kodeReseller, &kodeProduk, &jumlahTrx, &totalLaba, &kodeOperator, &NamaReseller, &sortOrder)
		if err != nil {
			log.Printf("failed to scan nested query row: %v", err)
			continue
		}

		switch tipe {
		case "TOTAL":
			result.TotalTransaksi = jumlahTrx
			result.TotalLaba = totalLaba

		case "RESELLER":
			reseller := &model.ResellerAnalytics{
				KodeReseller:    kodeReseller,
				JumlahTransaksi: jumlahTrx,
				TotalLaba:       totalLaba,
				Nama:            NamaReseller,
				ProdukBreakdown: []model.ProdukAnalytics{},
			}
			resellerMap[kodeReseller] = reseller
			resellerOrder = append(resellerOrder, kodeReseller)

		case "PRODUK":
			if reseller, exists := resellerMap[kodeReseller]; exists {
				produk := model.ProdukAnalytics{
					KodeProduk:      kodeProduk,
					JumlahTransaksi: jumlahTrx,
					TotalLaba:       totalLaba,
					KodeOperator:    kodeOperator,
				}
				reseller.ProdukBreakdown = append(reseller.ProdukBreakdown, produk)
			}
		}
	}

	if err = rows.Err(); err != nil {
		log.Printf("error iterating nested query rows: %v", err)
		return nil, fmt.Errorf("error iterating nested query rows: %w", err)
	}

	// Build final result maintaining order
	var finalResellers []model.ResellerAnalytics
	for _, kodeReseller := range resellerOrder {
		if reseller, exists := resellerMap[kodeReseller]; exists {
			finalResellers = append(finalResellers, *reseller)
		}
	}

	result.TransaksiPerReseller = finalResellers

	return &result, nil
}

func (repo *ProductRepository) TransaksiReseller(c context.Context, kodeReseller string, startDate, endDate string) ([]model.TransaksiDetails, error) {
	whereConditions := ""
	args := []interface{}{}
	if kodeReseller != "" {
		whereConditions += " AND kode_reseller LIKE @kodereseller"
		args = append(args, sql.Named("kodereseller", "%"+kodeReseller+"%"))
	}

	if startDate != "" {
		whereConditions += " AND CAST(tgl_status AS DATE) >= @startdate"
		args = append(args, sql.Named("startdate", startDate))
	}

	if startDate != "" {
		whereConditions += " AND CAST(tgl_status AS DATE) <= @enddate"
		args = append(args, sql.Named("enddate", startDate))
	}
	query := `
		SELECT 
			t.kode_produk,
			p.nama,
			t.tujuan,
			t.sn,
			t.harga_beli,
			t.harga
		FROM transaksi t
		LEFT JOIN produk p ON t.kode_produk = p.kode` + whereConditions + `
		ORDER BY t.kode_produk ASC
	`

	rows, err := repo.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var products []model.TransaksiDetails
	for rows.Next() {
		select {
		case <-c.Done():
			return nil, c.Err()
		default:
		}

		var product model.TransaksiDetails
		err := rows.Scan(
			&product.KodeProduk,
			&product.NamaProduk,
			&product.Tujuan,
			&product.Sn,
			&product.HargaBeli,
			&product.Harga,
		)
		if err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}
		products = append(products, product)
	}

	// Check untuk error dari rows.Next()
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration failed: %w", err)
	}

	return products, nil
}
func (repo *ProductRepository) Report(c context.Context, startDate string, endDate string, groupByTimeRange bool) ([]model.TopProductsBestSeller, error) {

	if startDate != "" {
		if _, err := time.Parse("2006-01-02", startDate); err != nil {
			return nil, fmt.Errorf("invalid startDate format: %w", err)
		}
	}
	if endDate != "" {
		if _, err := time.Parse("2006-01-02", endDate); err != nil {
			return nil, fmt.Errorf("invalid endDate format: %w", err)
		}
	}

	var query string
	var args []interface{}

	if groupByTimeRange {
		// Query dengan grouping berdasarkan jam
		query = `
		SELECT 
			CASE 
				WHEN DATEPART(HOUR, tgl_status) = 0 THEN '00:00 - 01:00'
				WHEN DATEPART(HOUR, tgl_status) = 1 THEN '01:00 - 02:00'
				WHEN DATEPART(HOUR, tgl_status) = 2 THEN '02:00 - 03:00'
				WHEN DATEPART(HOUR, tgl_status) = 3 THEN '03:00 - 04:00'
				WHEN DATEPART(HOUR, tgl_status) = 4 THEN '04:00 - 05:00'
				WHEN DATEPART(HOUR, tgl_status) = 5 THEN '05:00 - 06:00'
				WHEN DATEPART(HOUR, tgl_status) = 6 THEN '06:00 - 07:00'
				WHEN DATEPART(HOUR, tgl_status) = 7 THEN '07:00 - 08:00'
				WHEN DATEPART(HOUR, tgl_status) = 8 THEN '08:00 - 09:00'
				WHEN DATEPART(HOUR, tgl_status) = 9 THEN '09:00 - 10:00'
				WHEN DATEPART(HOUR, tgl_status) = 10 THEN '10:00 - 11:00'
				WHEN DATEPART(HOUR, tgl_status) = 11 THEN '11:00 - 12:00'
				WHEN DATEPART(HOUR, tgl_status) = 12 THEN '12:00 - 13:00'
				WHEN DATEPART(HOUR, tgl_status) = 13 THEN '13:00 - 14:00'
				WHEN DATEPART(HOUR, tgl_status) = 14 THEN '14:00 - 15:00'
				WHEN DATEPART(HOUR, tgl_status) = 15 THEN '15:00 - 16:00'
				WHEN DATEPART(HOUR, tgl_status) = 16 THEN '16:00 - 17:00'
				WHEN DATEPART(HOUR, tgl_status) = 17 THEN '17:00 - 18:00'
				WHEN DATEPART(HOUR, tgl_status) = 18 THEN '18:00 - 19:00'
				WHEN DATEPART(HOUR, tgl_status) = 19 THEN '19:00 - 20:00'
				WHEN DATEPART(HOUR, tgl_status) = 20 THEN '20:00 - 21:00'
				WHEN DATEPART(HOUR, tgl_status) = 21 THEN '21:00 - 22:00'
				WHEN DATEPART(HOUR, tgl_status) = 22 THEN '22:00 - 23:00'
				WHEN DATEPART(HOUR, tgl_status) = 23 THEN '23:00 - 00:00'
			END as time_range,
			DATEPART(HOUR, tgl_status) as hour_order,
			COUNT(*) as jumlah_transaksi,
			COALESCE(SUM(harga - harga_beli), 0) as total_laba,
			COUNT(DISTINCT kode_reseller) as count_member
		FROM transaksi
		WHERE status = 20
		`

		if startDate != "" {
			query += " AND CAST(tgl_status AS DATE) >= @startdate"
			args = append(args, sql.Named("startdate", startDate))
		}

		if endDate != "" {
			query += " AND CAST(tgl_status AS DATE) <= @enddate"
			args = append(args, sql.Named("enddate", endDate))
		}

		query += `
		GROUP BY DATEPART(HOUR, tgl_status)
		ORDER BY DATEPART(HOUR, tgl_status)
		`

	} else {
		query = `
		SELECT 
			'Total' as time_range,
			0 as hour_order,
			COUNT(*) as jumlah_transaksi,
			COALESCE(SUM(harga - harga_beli), 0) as total_laba,
			COUNT(DISTINCT kode_reseller) as count_member
		FROM transaksi
		WHERE status = 20
		`

		if startDate != "" {
			query += " AND CAST(tgl_status AS DATE) >= @startdate"
			args = append(args, sql.Named("startdate", startDate))
		}

		if endDate != "" {
			query += " AND CAST(tgl_status AS DATE) <= @enddate"
			args = append(args, sql.Named("enddate", endDate))
		}
	}

	// Execute query
	rows, err := repo.db.QueryContext(c, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query execution failed: %w", err)
	}
	defer rows.Close()

	// Scan results
	var products []model.TopProductsBestSeller
	for rows.Next() {
		// Check context cancellation
		select {
		case <-c.Done():
			return nil, c.Err()
		default:
		}

		var product model.TopProductsBestSeller
		var hourOrder int

		err := rows.Scan(
			&product.TimeRange,
			&hourOrder,
			&product.CountTrx,
			&product.TotalProfit,
			&product.CountMember,
		)
		if err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}

		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration failed: %w", err)
	}

	return products, nil
}
