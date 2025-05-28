package entity

import "time"

type User struct {
	ID           int64     `json:"id"`
	Username     string    `json:"username"`
	Role         string    `json:"role"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}

type (
	AccessClaims struct {
		UserID int64  `mapstructure:"user_id"`
		Role   string `mapstructure:"role"`
		Exp    int64  `mapstructure:"exp"`
		Iat    int64  `mapstructure:"iat"`
	}

	RefreshClaims struct {
		UserID int64 `mapstructure:"user_id"`
		Exp    int64 `mapstructure:"exp"`
		Iat    int64 `mapstructure:"iat"`
	}
)
