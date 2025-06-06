package controller

import (
	"fmt"
	"net/http"

	"github.com/Van-programan/Forum_GO/internal/controller/request"
	_ "github.com/Van-programan/Forum_GO/internal/controller/response"
	_ "github.com/Van-programan/Forum_GO/internal/entity"
	"github.com/Van-programan/Forum_GO/internal/usecase"
	"github.com/gin-gonic/gin"

	"github.com/rs/zerolog"
)

type AuthHandler struct {
	Usecase usecase.AuthUsecase
	Log     *zerolog.Logger
}

const (
	registerOp     = "AuthHandler.Register"
	loginOp        = "AuthHandler.Login"
	refreshOp      = "AuthHandler.Refresh"
	logoutOp       = "AuthHandler.Logout"
	checkSessionOp = "AuthHandler.CheckSession"
)

// Register godoc
// @Summary Register a new user
// @Description Creates a new user account and returns user information along with an access token. A refresh token is set as an HTTP-only cookie.
// @Tags auth
// @Accept json
// @Produce json
// @Param user body request.RegisterRequest true "User Credentials"
// @Success 200 {object} response.RegisterSuccessResponse "Successfully registered"
// @Failure 400 {object} response.ErrorResponseAuth "Invalid request payload"
// @Failure 500 {object} response.ErrorResponseAuth "Internal server error"
// @Router /register [post]
func (ah *AuthHandler) Register(c *gin.Context) {
	log := ah.getRequestLogger(c).With().Str("op", registerOp).Logger()

	var req request.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn().Err(err).Msg("Failed to bind request")
		fmt.Println("ошибка")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := ah.Usecase.Register(c.Request.Context(), req.Username, req.Role, req.Password)
	if err != nil {
		log.Error().Err(err).Msg("Failed to register user")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie("refresh_token", res.Tokens.RefreshToken, 3600*24*30, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"user": res.User, "access_token": res.Tokens.AccessToken})
}

// Login godoc
// @Summary Log in an existing user
// @Description Authenticates a user and returns user information along with a new access token. A new refresh token is set as an HTTP-only cookie, and any existing refresh token in the cookie is invalidated.
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body request.LoginRequest true "User Login Credentials"
// @Success 200 {object} response.LoginSuccessResponse "Successfully logged in"
// @Failure 400 {object} response.ErrorResponseAuth "Invalid request payload"
// @Failure 401 {object} response.ErrorResponseAuth "Invalid credentials"
// @Router /login [post]
func (ah *AuthHandler) Login(c *gin.Context) {
	log := ah.getRequestLogger(c).With().Str("op", loginOp).Logger()

	var req request.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn().Err(err).Msg("Failed to bind request")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		fmt.Println("ошибка")
		return
	}

	oldRefreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		if err == http.ErrNoCookie {
			oldRefreshToken = ""
		} else {
			log.Warn().Err(err).Msg("Error reading refresh_token cookie")
			oldRefreshToken = ""
		}
	}

	res, err := ah.Usecase.Login(c.Request.Context(), req.Username, req.Password, oldRefreshToken)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to login user")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	c.SetCookie("refresh_token", res.Tokens.RefreshToken, 3600*24*30, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"user": res.User, "access_token": res.Tokens.AccessToken})
}

// Refresh godoc
// @Summary Refresh access token
// @Description Uses a refresh token to generate a new access token and a new refresh token. Refresh token is set as an HTTP-only cookie.
// @Tags auth
// @Produce json
// @Success 200 {object} response.RefreshSuccessResponse "Successfully refreshed tokens"
// @Failure 401 {object} response.ErrorResponseAuth "Refresh token required or invalid/expired refresh token"
// @Router /refresh [post]
func (ah *AuthHandler) Refresh(c *gin.Context) {
	log := ah.getRequestLogger(c).With().Str("op", refreshOp).Logger()

	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		log.Warn().Err(err).Msg("Failed to get refresh token from cookie")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token required"})
		return
	}

	tokens, err := ah.Usecase.Refresh(c.Request.Context(), refreshToken)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to refresh token")
		c.SetCookie("refresh_token", "", -1, "/", "", false, true)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired refresh token"})
		return
	}

	c.SetCookie("refresh_token", tokens.Tokens.RefreshToken, 3600*24*30, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"access_token": tokens.Tokens.AccessToken})
}

// Logout godoc
// @Summary Log out a user
// @Description Logs out a user by deleting the refresh token from the server and clearing the refresh token cookie.
// @Tags auth
// @Produce json
// @Success 200 {object} response.LogoutSuccessResponse "Successfully logged out"
// @Router /logout [post]
func (ah *AuthHandler) Logout(c *gin.Context) {
	log := ah.getRequestLogger(c).With().Str("op", logoutOp).Logger()
	refreshToken, err := c.Cookie("refresh_token")
	if err == nil && refreshToken != "" {
		_ = ah.Usecase.Logout(c.Request.Context(), refreshToken)
	} else {
		log.Error().Err(err).Msg("Failed to get refresh token from cookie")
	}

	c.SetCookie("refresh_token", "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
}

func (ah *AuthHandler) getRequestLogger(c *gin.Context) *zerolog.Logger {
	reqLog := ah.Log.With().
		Str("method", c.Request.Method).
		Str("path", c.Request.URL.Path).
		Str("remote_addr", c.ClientIP())

	logger := reqLog.Logger()
	return &logger
}

func (ah *AuthHandler) CheckSession(c *gin.Context) {
	log := ah.getRequestLogger(c).With().Str("op", checkSessionOp).Logger()

	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		log.Warn().Err(err).Msg("Failed to get refresh token from cookie")
		c.JSON(http.StatusUnauthorized, gin.H{"user": nil, "is_active": false})
		return
	}

	res, err := ah.Usecase.IsSessionActive(c.Request.Context(), refreshToken)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to check session")
		c.JSON(http.StatusUnauthorized, gin.H{"user": nil, "is_active": false})
		return
	}
	fmt.Println(res.User, res.IsActive)
	c.JSON(http.StatusOK, gin.H{"user": res.User, "is_active": res.IsActive})
}
