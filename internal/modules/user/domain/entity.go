package domain

import "time"

type UserStatus string

const (
	UserStatusActive   UserStatus = "active"
	UserStatusDisabled UserStatus = "disabled"
)

type User struct {
	ID           string
	Email        string
	PasswordHash string
	NickName     string
	Status       UserStatus
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    time.Time
}