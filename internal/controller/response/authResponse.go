package response

import "github.com/Van-programan/Forum_GO/internal/entity"

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RegisterResponse struct {
	User   entity.User `json:"user"`
	Tokens Tokens      `json:"tokens"`
}

type RegisterSuccessResponse struct {
	User        entity.User `json:"user"`
	AccessToken string      `json:"access_token"`
}

type LoginResponse struct {
	User   entity.User `json:"user"`
	Tokens Tokens      `json:"tokens"`
}

type LoginSuccessResponse struct {
	User        entity.User `json:"user"`
	AccessToken string      `json:"access_token"`
}

type RefreshResponse struct {
	Tokens Tokens `json:"tokens"`
}

type RefreshSuccessResponse struct {
	AccessToken string `json:"access_token"`
}

type LogoutSuccessResponse struct {
	Message string `json:"message" example:"logged out successfully"`
}

type IsSessionActiveResponse struct {
	User     *entity.User `json:"user"`
	IsActive bool         `json:"is_active"`
}

type ErrorResponseAuth struct {
	Error string `json:"error" example:"error message"`
}
