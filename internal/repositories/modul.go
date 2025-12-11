package repositories

import (
	"context"
	"database/sql"
)

type ModulType struct {
	Code  string   `json:"kode"`
	Label string   `json:"label"`
	Saldo *float64 `json:"saldo"`
}

type ModulOtomax struct {
	Db *sql.DB
}

func NewModulOtomax(DB *sql.DB) *ModulOtomax {
	return &ModulOtomax{
		Db: DB,
	}
}
func (m *ModulOtomax) GetAllModulOtomax(c context.Context) ([]ModulType, error) {
	query := `
		SELECT 
				kode,
				label,
				saldo
			FROM modul m1
			WHERE aktif = 1
			AND kode = (
				SELECT MIN(m2.kode)
				FROM modul m2
				WHERE m2.tujuan = m1.tujuan
					AND m2.aktif = 1
			)
			ORDER BY kode ASC
		`

	var datas []ModulType
	rows, err := m.Db.QueryContext(c, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var kode, label string
		var saldo *float64
		if err := rows.Scan(&kode, &label, &saldo); err != nil {
			return nil, err
		}
		data := ModulType{
			Code:  kode,
			Label: label,
			Saldo: saldo,
		}
		datas = append(datas, data)
	}
	return datas, nil
}
