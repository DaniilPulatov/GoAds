package rest

import (
	"net/http"

	authHandle "ads-service/internal/rest/handlers/auth"
	userHandle "ads-service/internal/rest/handlers/user"
	"ads-service/internal/rest/middleware"

	"github.com/gin-gonic/gin"
)

type Server struct {
	mux         *gin.Engine
	authHandler *authHandle.AuthHandler
	userHandler *userHandle.UserHandler
	mv          *middleware.Middleware
}

func NewServer(mux *gin.Engine, authHandler *authHandle.AuthHandler, userHandler *userHandle.UserHandler,
	mdvr *middleware.Middleware) *Server {
	mux.Use(gin.Recovery())
	mux.Use(gin.Logger())
	mux.LoadHTMLGlob("web/*")

	return &Server{
		mux:         mux,
		authHandler: authHandler,
		userHandler: userHandler,
		mv:          mdvr,
	}

}

func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	s.mux.ServeHTTP(writer, request)
}

func (s *Server) Init() {
	baseGroup := s.mux.Group("/api/v1")
	{
		authGroup := baseGroup.Group("/auth")
		{
			authGroup.POST("/register", s.authHandler.Register)
			authGroup.POST("/login", s.authHandler.Login)
		}
	}

}
