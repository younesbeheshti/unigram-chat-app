package models

type User struct {
	ID int64	`json:"id"`
	Name string	`json:"name"`
	Password string	`json:"password"`
}

type Respnse struct {
	UserID int64 `json:"userid,omitempty"`
	Message string `json:"message,omitempty"`
}