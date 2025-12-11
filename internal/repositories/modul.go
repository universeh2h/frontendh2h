package repositories

import (
	"context"
	"database/sql"
)

type ModulType struct {
	Code     string   `json:"kode"`
	Label    string   `json:"label"`
	Saldo    *float64 `json:"saldo"`
	TotalTrx int      `json:"total_trx"`
}

type ModulOtomax struct {
	Db *sql.DB
}

func NewModulOtomax(DB *sql.DB) *ModulOtomax {
	return &ModulOtomax{
		Db: DB,
	}
}
func (m *ModulOtomax) GetAllModulOtomax(c context.Context, tglStatus string) ([]ModulType, error) {
	query := `
		SELECT 
			m1.kode,
			m1.label,
			m1.saldo,
			COUNT(t.kode) as total_trx
		FROM modul m1
		LEFT JOIN transaksi t 
			ON m1.kode = t.kode_modul
			AND CAST(t.tgl_status AS DATE) = @p1
		WHERE m1.aktif = 1
			AND m1.kode = (
				SELECT MIN(m2.kode)
				FROM modul m2
				WHERE m2.tujuan = m1.tujuan
					AND m2.aktif = 1
			)
		GROUP BY m1.kode, m1.label, m1.saldo
		ORDER BY m1.kode ASC
	`

	var datas []ModulType
	rows, err := m.Db.QueryContext(c, query, tglStatus)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var kode, label string
		var saldo *float64
		var totalTrx int
		if err := rows.Scan(&kode, &label, &saldo, &totalTrx); err != nil {
			return nil, err
		}
		data := ModulType{
			Code:     kode,
			Label:    label,
			Saldo:    saldo,
			TotalTrx: totalTrx,
		}
		datas = append(datas, data)
	}
	return datas, nil
}
