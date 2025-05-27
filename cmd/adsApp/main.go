package main

import (
	"ads-service/internal/migrations"
	adRepository "ads-service/internal/repository/ad"
	adFileRepository "ads-service/internal/repository/adFile"
	authRepository "ads-service/internal/repository/auth"
	userRepository "ads-service/internal/repository/user"
	authHandler "ads-service/internal/rest/handlers/auth"
	mv "ads-service/internal/rest/middleware"
	adminService "ads-service/internal/usecase/admin"
	authService "ads-service/internal/usecase/auth"
	userService "ads-service/internal/usecase/user"
	"time"

	"ads-service/internal/rest"
	"ads-service/pkg/db"
	"errors"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"go.uber.org/dig"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("Error loading .env file, using default values")
		os.Exit(1)
	}
	log.Println("Env vars loaded successfully")
	var (
		port = os.Getenv("PORT")
		host = os.Getenv("HOST")
		dsn  = os.Getenv("DATABASE_URL")
	)

	if err := migrations.New("file://internal/migrations", dsn); err != nil {
		log.Println("Error running migrations:", err)
		os.Exit(1)
	}

	if err := execute(host, port, dsn); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func execute(host, port, dsn string) error {
	deps := []interface{}{
		func() (*pgxpool.Pool, error) {
			return db.NewDB(dsn)
		},
		func() *gin.Engine {
			return gin.New()
		},
		authHandler.NewAuthHandler,

		authService.NewAuthService,
		adminService.NewAdminService,
		userService.NewUserService,

		authRepository.NewAuthRepo,
		userRepository.NewUserRepo,
		adFileRepository.NewAdFileRepo,
		adRepository.NewAdRepo,

		mv.NewMiddleware,

		http.NewServeMux,
		rest.NewServer,
		func(server *rest.Server) *http.Server {
			return &http.Server{
				Addr:              net.JoinHostPort(host, port),
				Handler:           server,
				ReadHeaderTimeout: 10 * time.Second,
			}
		},
	}

	container := dig.New()
	for _, dep := range deps {
		if err := container.Provide(dep); err != nil {
			log.Println("failed to provide dependency:", err)
			return errors.New("failed to provide dependency")
		}
	}

	err := container.Invoke(func(server *rest.Server) {
		server.Init()
	})
	if err != nil {
		log.Panicln("failed to initialize server:", err)
		return errors.New("failed to initialize server")
	}

	err = container.Invoke(func(server *http.Server) error {
		return server.ListenAndServe()
	})

	if err != nil {
		log.Println("failed to start server:", err)
		return errors.New("failed to start server")
	}
	return nil
}
