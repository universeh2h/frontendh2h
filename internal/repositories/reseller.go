package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"log"
)

type ProductResponse struct {
	KodeProduk string  `json:"kode_produk"`
	JumlahTRX  int     `json:"count_trx"`
	TotalLaba  float64 `json:"total_laba"`
}

type TrxTerCuan struct {
	KodeProduk string  `json:"kode_produk"`
	JumlahTRX  int     `json:"count_trx"`
	TotalLaba  float64 `json:"total_laba"`
	Rata2Laba  float64 `json:"rata_laba"`
}

func (repo *ProductRepository) GetProductTrxTerbanyak(c context.Context, startDate string, endDate string, kodeReseller string) ([]ProductResponse, error) {
	whereConditions := ""
	args := []interface{}{}

	if startDate != "" {
		whereConditions += " AND CAST(tgl_status AS DATE) >= @startdate"
		args = append(args, sql.Named("startdate", startDate))
	}
	if kodeReseller != "" {
		whereConditions += " AND kode_reseller=@kodeReseller"
		args = append(args, sql.Named("kodeReseller", kodeReseller))
	}

	if endDate != "" {
		whereConditions += " AND CAST(tgl_status AS DATE) <= @enddate"
		args = append(args, sql.Named("enddate", endDate))
	}

	query := `
		SELECT TOP 10
			kode_produk,
			COUNT(*) as jumlah_transaksi,
			SUM(harga - harga_beli) as total_laba
		FROM transaksi
		WHERE  status = 20` + whereConditions + `
		GROUP BY kode_produk
		ORDER BY jumlah_transaksi DESC;
	`
	rows, err := repo.db.QueryContext(c, query, args...)
	if err != nil {
		log.Printf("failed to execute nested analytics query: %v", err)
		return nil, fmt.Errorf("failed to execute nested analytics query: %w", err)
	}
	defer rows.Close()

	var results []ProductResponse
	for rows.Next() {
		var product ProductResponse
		err := rows.Scan(
			&product.KodeProduk,
			&product.JumlahTRX,
			&product.TotalLaba,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transaksi: %w", err)
		}

		results = append(results, product)
	}

	// Check for iteration errors
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during row iteration: %w", err)
	}

	return results, nil

}

func (repo *ProductRepository) GetTrxTercuan(c context.Context, startDate string, endDate string, kodeReseller string) ([]TrxTerCuan, error) {
	whereConditions := ""
	args := []interface{}{}

	if startDate != "" {
		whereConditions += " AND CAST(tgl_status AS DATE) >= @startdate"
		args = append(args, sql.Named("startdate", startDate))
	}
	if kodeReseller != "" {
		whereConditions += " AND kode_reseller=@kodeReseller"
		args = append(args, sql.Named("kodeReseller", kodeReseller))
	}

	if endDate != "" {
		whereConditions += " AND CAST(tgl_status AS DATE) <= @enddate"
		args = append(args, sql.Named("enddate", endDate))
	}

	query := `
	SELECT TOP 10
		kode_produk,
		COUNT(*) as jumlah_transaksi,
		SUM(harga - harga_beli) as total_laba,
		AVG(harga - harga_beli) as rata_rata_laba
	FROM transaksi
	WHERE  status = 20` + whereConditions + `
	GROUP BY kode_produk
	ORDER BY total_laba DESC;
	`
	rows, err := repo.db.QueryContext(c, query, args...)
	if err != nil {
		log.Printf("failed to execute nested analytics query: %v", err)
		return nil, fmt.Errorf("failed to execute nested analytics query: %w", err)
	}
	defer rows.Close()

	var results []TrxTerCuan
	for rows.Next() {
		var product TrxTerCuan
		err := rows.Scan(
			&product.KodeProduk,
			&product.JumlahTRX,
			&product.TotalLaba,
			&product.Rata2Laba,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transaksi: %w", err)
		}

		results = append(results, product)
	}

	// Check for iteration errors
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during row iteration: %w", err)
	}

	return results, nil

}
