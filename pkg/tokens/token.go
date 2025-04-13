package tokens

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const (
	AccessTokenTTL  = 15 * time.Minute
	RefreshTokenTTL = 35 * 24 * time.Hour
	BcryptCost      = 12
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrTokenExpired = errors.New("token expired")
)

type TokenManager interface {
	GenerateAccessToken(userID int64) (string, error)
	GenerateRefreshToken() (string, error)
	ParseAccessToken(token string) (int64, error)
	HashPassword(password string) (string, error)
	CheckPasswordHash(password, hash string) bool
}

type tokenManager struct {
	signingKey string
}

func NewTokenManager(signingKey string) TokenManager {
	return &tokenManager{signingKey: signingKey}
}

func (tm *tokenManager) GenerateAccessToken(userID int64) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   string(rune(userID)),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(AccessTokenTTL)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(tm.signingKey))
}

func (tm *tokenManager) GenerateRefreshToken() (string, error) {
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(RefreshTokenTTL)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(tm.signingKey))
}

func (tm *tokenManager) ParseAccessToken(tokenString string) (int64, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(tm.signingKey), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return 0, ErrTokenExpired
		}
		return 0, ErrInvalidToken
	}

	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		userID := []byte(claims.Subject)[0]
		return int64(userID), nil
	}

	return 0, ErrInvalidToken
}

func (tm *tokenManager) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), BcryptCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (tm *tokenManager) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
