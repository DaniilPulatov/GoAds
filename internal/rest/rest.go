// @title           My Ads API
// @version         1.0
// @description     REST API for ads project

// @contact.name   Support
// @contact.email  d.pulatov@student.inha.uz, s.raxmatov@student.inha.uz

// @host      localhost:8080
// @BasePath  /api/v1

// @schemes http

package rest

import (
	"ads-service/internal/rest/handlers/admin"
	"ads-service/internal/rest/handlers/user"
	"net/http"

	authHandle "ads-service/internal/rest/handlers/auth"
	"ads-service/internal/rest/middleware"

	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	mux          *gin.Engine
	authHandler  *authHandle.AuthHandler
	adminHandler *admin.AdminHandler
	userHandler  *user.UserHandler
	mv           *middleware.Middleware
}

func NewServer(mux *gin.Engine, authHandler *authHandle.AuthHandler, mv *middleware.Middleware,
	adminHandler *admin.AdminHandler, userHandler *user.UserHandler) *Server {
	mux.Use(gin.Recovery())
	mux.Use(gin.Logger())

	return &Server{
		mux:          mux,
		authHandler:  authHandler,
		adminHandler: adminHandler,
		userHandler:  userHandler,
		mv:           mv,
	}

}

func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	s.mux.ServeHTTP(writer, request)
}

func (s *Server) Init() {
	s.mux.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	const basePath = "/api/v1"
	baseGroup := s.mux.Group(basePath)

	// Маршруты для аутентификации
	authGroup := baseGroup.Group("/auth")
	authGroup.POST("/register", s.authHandler.Register)
	authGroup.POST("/login", s.authHandler.Login)

	// Пользовательские маршруты
	userGroup := baseGroup.Group("/ads")
	userGroup.Use(s.mv.UserAuth())
	userGroup.POST("/create", s.userHandler.CreateDraft)
	userGroup.GET("/my", s.userHandler.GetMyAds)
	userGroup.PUT("/:id", s.userHandler.UpdateMyAd)
	userGroup.DELETE("/:id", s.userHandler.DeleteMyAd)
	userGroup.POST("/:id/submit", s.userHandler.SubmitForModeration)
	userGroup.POST("/:id/image", s.userHandler.AddImageToMyAd)
	userGroup.GET("/:id/image", s.userHandler.GetImagesToMyAd)
	userGroup.DELETE("/:id/image/:fid", s.userHandler.DeleteMyAdImage)
	userGroup.GET("/filter", s.userHandler.GetMyAdsByFilter)

	// Админские маршруты
	adminGroup := baseGroup.Group("/admin")
	adminGroup.Use(s.mv.AdminAuth())
	adminGroup.GET("/ads", s.adminHandler.GetAllAds)
	adminGroup.GET("/stats", s.adminHandler.GetStatistics)
	adminGroup.DELETE("/ads/:id", s.adminHandler.DeleteAd)
	adminGroup.POST("/ads/:id/approve", s.adminHandler.Approve)
	adminGroup.POST("/ads/:id/reject", s.adminHandler.Reject)
}
