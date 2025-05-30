package route

import (
	"time"

	"github.com/Van-programan/Forum_GO/internal/client"
	"github.com/Van-programan/Forum_GO/internal/controller"
	"github.com/Van-programan/Forum_GO/internal/controller/middleware"
	"github.com/Van-programan/Forum_GO/internal/usecase"
	"github.com/Van-programan/Forum_GO/internal/ws"
	"github.com/Van-programan/Forum_GO/pkg/jwt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/Van-programan/Forum_GO/docs" // swagger docs
)

func NewForumRouter(engine *gin.Engine, categoryUsecase usecase.CategoryUsecase,
	topicUsecase usecase.TopicUsecase, postUsecase usecase.PostUsecase,
	jwt *jwt.JWT, log *zerolog.Logger, hub *ws.Hub, chatUsecase usecase.ChatUsecase,
	userClient client.UserClient) {

	// Initialize Swagger info
	docs.SwaggerInfo.Title = "Forum Service API"
	docs.SwaggerInfo.Description = "API for forum service"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:3101"
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"http"}

	categoryHandler := &controller.CategoryHandler{
		Usecase: categoryUsecase,
		Log:     log,
	}
	topicHandler := &controller.TopicHandler{
		Usecase: topicUsecase,
		Log:     log,
	}
	postHandler := &controller.PostHandler{Usecase: postUsecase}
	auth := middleware.NewAuthMiddleware(jwt)
	chatHandler := controller.NewChatHandler(hub, chatUsecase, userClient, log)

	engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	engine.GET("/ws", auth.ChatAuth(), chatHandler.ServeWs)

	categories := engine.Group("/categories")
	{
		categories.GET("", categoryHandler.GetAll)
		categories.GET("/:id", categoryHandler.GetByID)

		adminCategories := categories.Group("")
		adminCategories.Use(auth.Auth(), middleware.RequireAdmin())
		{
			adminCategories.POST("", categoryHandler.Create)
			adminCategories.DELETE("/:id", categoryHandler.Delete)
			adminCategories.PATCH("/:id", categoryHandler.Update)
		}
	}

	engine.GET("/categories/topics/:id", topicHandler.GetByCategory)
	engine.POST("/categories/topics/:id/", auth.Auth(), topicHandler.Create)

	engine.GET("/topics/:id", topicHandler.GetByID)
	topics := engine.Group("/topics").Use(auth.Auth())
	{
		topics.DELETE("/:id", topicHandler.Delete)
		topics.PATCH("/:id", topicHandler.Update)
	}

	engine.GET("/topics/:id/posts", postHandler.GetByTopic)
	engine.POST("/topics/:id/posts", auth.Auth(), postHandler.Create)

	posts := engine.Group("/posts").Use(auth.Auth())
	{
		posts.DELETE("/:id", postHandler.Delete)
		posts.PATCH("/:id", postHandler.Update)
	}

	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
