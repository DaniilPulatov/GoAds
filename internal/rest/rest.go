package rest

import (
	"net/http"

	"ads-service/internal/rest/handlers/auth"

	"github.com/gin-gonic/gin"
)

type Server struct {
	mux         *gin.Engine
	authHandler *auth.AuthHandler
	// mv          *auth.Midlleware
}

func NewServer(mux *gin.Engine, authHandler *auth.AuthHandler) *Server {
	mux.Use(gin.Recovery())
	mux.Use(gin.Logger())

	return &Server{
		mux:         mux,
		authHandler: authHandler,
	}
}

func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	s.mux.ServeHTTP(writer, request)
}

func (s *Server) Init() {
	baseGroup := s.mux.Group("/ads/api/v1")
	{
		authGroup := baseGroup.Group("/auth")
		{
			authGroup.POST("/register", s.authHandler.Register)
			authGroup.POST("/login", s.authHandler.Login)
		}
	}

}
