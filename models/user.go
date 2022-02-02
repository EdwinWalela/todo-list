package models

type User struct {
	Id       int64  `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Is_admin int8   `json:"is_admin"`
}
