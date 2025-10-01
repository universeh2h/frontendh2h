package model

// Struct untuk hasil nested analytics
type AnalyticsResult struct {
	TotalTransaksi       int                 `json:"total_transaksi"`
	TotalLaba            float64             `json:"total_laba"`
	TransaksiPerReseller []ResellerAnalytics `json:"transaksi_per_reseller"`
}

type ResellerAnalytics struct {
	KodeReseller    string            `json:"kode_reseller"`
	JumlahTransaksi int               `json:"jumlah_transaksi"`
	TotalLaba       float64           `json:"total_laba"`
	Nama            string            `json:"namaReseller"`
	ProdukBreakdown []ProdukAnalytics `json:"produk_breakdown"`
}

type ProdukAnalytics struct {
	KodeProduk      string  `json:"kode_produk"`
	JumlahTransaksi int     `json:"jumlah_transaksi"`
	KodeOperator    string  `json:"kodeOperator"`
	TotalLaba       float64 `json:"total_laba"`
}

type PaginationParams struct {
	Search    *string `json:"search,omitempty"`
	StartDate string  `json:"startDate"`
	EndDate   string  `json:"endDate"`
}

type TransaksiDetails struct {
	KodeProduk string  `json:"kodeProduk"`
	NamaProduk *string `json:"namaProduk"`
	Tujuan     string  `json:"tujuan"`
	Sn         *string `json:"sn"`
	Harga      float64 `json:"hargaJual"`
	HargaBeli  float64 `json:"hargaBeli"`
}
