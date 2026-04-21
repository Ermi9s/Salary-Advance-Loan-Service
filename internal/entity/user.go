package entity

import "time"

type UserRole string

const (
	Admin    UserRole = "admin"
	Uploader UserRole = "uploader"
)

type User struct {
	ID           string    `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	Role         UserRole  `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
}

type AuthUser struct {
	Username string `json:"username" binding:"required,min=3"`
	Password string `json:"password" binding:"required,min=8"`
}

type TokenClaims struct {
	Username string   `json:"username"`
	Role     UserRole `json:"role"`
}
