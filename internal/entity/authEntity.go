package entity

import "time"

type User struct {
	ID           int64     `json:"id" db:"id"`
	Username     string    `json:"username" db:"username"`
	Email        string    `json:"email" db:"email"`
	Password     string    `json:"-" db:"password_hash"`
	RegisteredAt time.Time `json:"registered_at" db:"registered_at"`
}

type Session struct {
	ID                    int64     `json:"id" db:"id"`
	UserID                int64     `json:"user_id" db:"user_id"`
	AccessToken           string    `json:"access_token" db:"access_token"`
	RefreshToken          string    `json:"refresh_token" db:"refresh_token"`
	ExpiresAtAccessToken  time.Time `json:"access_token_expires_at" db:"access_token_expires_at"`
	ExpiresAtRefreshToken time.Time `json:"refresh_token_expires_at" db:"refresh_token_expires_at"`
}
