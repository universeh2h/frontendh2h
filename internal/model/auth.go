package model

type User struct {
	Username string `json:"kode"`
}

type Login struct {
	Username string `json:"kode"`
	Password string `json:"password"`
}
