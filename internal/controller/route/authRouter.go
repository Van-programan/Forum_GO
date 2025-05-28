package route

import (
	"time"

	"github.com/Van-programan/Forum_GO/internal/controller"
	"github.com/Van-programan/Forum_GO/internal/usecase"
	"github.com/Van-programan/Forum_GO/pkg/jwt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewAuthRouter(engine *gin.Engine, usecase usecase.AuthUsecase, jwt *jwt.JWT, log *zerolog.Logger) {
	h := &controller.AuthHandler{
		Usecase: usecase,
		Log:     log,
	}

	engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"POST"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	engine.POST("/register", h.Register)
	engine.POST("/login", h.Login)
	engine.POST("/refresh", h.Refresh)
	engine.POST("/logout", h.Logout)
	engine.GET("/check-session", h.CheckSession)

	engine.GET("/swagger/auth/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
