package rest

import (
	"ads-service/internal/rest/handlers/admin"
	"ads-service/internal/rest/handlers/user"
	"net/http"

	authHandle "ads-service/internal/rest/handlers/auth"
	userHandle "ads-service/internal/rest/handlers/user"
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
	mux.LoadHTMLGlob("web/*")

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

	const (
		basePath = "/api/v1"
	)
	baseGroup := s.mux.Group(basePath)
	baseGroup.Use(s.middleware.UserAuth())
	{
		baseGroup.POST("/ad", s.userHandler.CreateDraft)
	}

	authGroup := baseGroup.Group("/auth")
	{
		authGroup.POST("/register", s.authHandler.Register)
		authGroup.POST("/login", s.authHandler.Login)
	}
}
