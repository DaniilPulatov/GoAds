package middlware

import (
	"ads-service/internal/usecase/auth"
	"ads-service/internal/usecase/user"
)

type Middleware struct {
	authService auth.AuthService
	userService user.UserAdvertisementService
}

func NewMiddleware(authService auth.AuthService, userService user.UserAdvertisementService) *Middleware {
	return &Middleware{
		authService: authService,
		userService: userService,
	}
}
