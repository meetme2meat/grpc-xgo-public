package model

type User struct {
	Id           int    `json:"id" gorm:"primaryKey;autoIncrement;not null"`
	Username     string `json:"username"`
	PasswordHash string `json:"password_hash"`
}
