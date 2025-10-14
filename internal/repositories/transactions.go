package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type TransactionData struct {
	TglEntri     time.Time `json:"tgl_entri"`
	TglStatus    time.Time `json:"tgl_status"`
	Sn           *string   `json:"sn"`
	Status       string    `json:"status"`
	TrxID        int64     `json:"trx_id"`
	ProductCode  string    `json:"product_code"`
	Tujuan       string    `json:"tujuan"`
	SellingPrice float64   `json:"selling_price"`
}

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{
		db: db,
	}
}

func (repo *TransactionRepository) GetTransactions(c context.Context) ([]TransactionData, error) {

	query := `
		SELECT 
			TOP 10
			tgl_entri,
			status,
			sn,
			tgl_status,
			kode,
			kode_produk,
			tujuan,
			harga
		FROM transaksi
		ORDER BY tgl_entri DESC
	`

	rows, err := repo.db.QueryContext(c, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions : %s", err.Error())
	}

	defer rows.Close()
	var transactions []TransactionData
	for rows.Next() {
		var transaksi TransactionData
		var statusInt int
		err := rows.Scan(
			&transaksi.TglEntri,
			&statusInt,
			&transaksi.Sn,
			&transaksi.TglStatus,
			&transaksi.TrxID,
			&transaksi.ProductCode,
			&transaksi.Tujuan,
			&transaksi.SellingPrice,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to get transactions : %s", err.Error())

		}

		transaksi.Status = VerifiedStatus(statusInt)

		transactions = append(transactions, transaksi)
	}

	return transactions, nil
}

func VerifiedStatus(status int) string {
	var (
		statusPendingString = "pending"
		statusSuccesString  = "success"
		statusFailedString  = "failed"
	)
	switch status {
	case 1, 2:
		return statusPendingString
	case 20:
		return statusSuccesString
	default:
		return statusFailedString
	}

}
