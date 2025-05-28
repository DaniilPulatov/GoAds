package rest

import (
	"ads-service/internal/rest/handlers/admin"
	"ads-service/internal/rest/handlers/user"
	"net/http"

	authHandle "ads-service/internal/rest/handlers/auth"
	"ads-service/internal/rest/middleware"

	"github.com/gin-gonic/gin"
)

type Server struct {
	mux          *gin.Engine
	authHandler  *authHandle.AuthHandler
	adminHandler *admin.AdminHandler
	userHandler  *user.UserHandler
	middleware   *middleware.Middleware
}

func NewServer(mux *gin.Engine, authHandler *authHandle.AuthHandler, middleware *middleware.Middleware,
	adminHandler *admin.AdminHandler, userHandler *user.UserHandler) *Server {
	mux.Use(gin.Recovery())
	mux.Use(gin.Logger())

	return &Server{
		mux:          mux,
		authHandler:  authHandler,
		adminHandler: adminHandler,
		userHandler:  userHandler,
		middleware:   middleware,
	}

}

func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	s.mux.ServeHTTP(writer, request)
}

func (s *Server) Init() {
	const basePath = "/api/v1"
	baseGroup := s.mux.Group(basePath)

	authGroup := baseGroup.Group("/auth")
	{
		authGroup.POST("/register", s.authHandler.Register)
		authGroup.POST("/login", s.authHandler.Login)
	}

	// Пользовательские маршруты
	userGroup := baseGroup.Group("/ads")
	userGroup.Use(s.middleware.UserAuth())
	{
		userGroup.POST("/create", s.userHandler.CreateDraft)
		userGroup.GET("/my", s.userHandler.GetMyAds)
		userGroup.PUT("/:id", s.userHandler.UpdateMyAd)
		userGroup.DELETE("/:id", s.userHandler.DeleteMyAd)
		userGroup.POST("/:id/submit", s.userHandler.SubmitForModeration)
		userGroup.POST("/:id/image", s.userHandler.AddImageToMyAd)
		userGroup.GET("/:id/image", s.userHandler.GetImagesToMyAd)
		userGroup.DELETE("/:id/image/:fid", s.userHandler.DeleteMyAdImage)
	}

	// Админские маршруты
	adminGroup := baseGroup.Group("/admin")
	adminGroup.Use(s.middleware.AdminAuth())
	{
		adminGroup.GET("/ads", s.adminHandler.GetAllAds)
		adminGroup.GET("/stats", s.adminHandler.GetStatistics)
		adminGroup.DELETE("/ads/:id", s.adminHandler.DeleteAd)
		adminGroup.POST("/ads/:id/approve", s.adminHandler.Approve)
		adminGroup.POST("/ads/:id/reject", s.adminHandler.Reject)
		adminGroup.DELETE("/ads/:id/image/:fid", s.adminHandler.DeleteImage)
	}
}
