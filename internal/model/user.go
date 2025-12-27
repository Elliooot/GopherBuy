package model

type User struct {
	ID           uint64 `gorm:"primarykey"`
	Username     string
	Email        string
	PasswordHash string
	PhoneNumber  string
}
