package user

import "ads-service/internal/usecase/user"

type UserHandler struct {
	userService user.UserAdvertisementService
}

func NewUserHandler(userService user.UserAdvertisementService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

type UpdateAdRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	CategoryID  int    `json:"category_id"`
}
