package admin

import "ads-service/internal/usecase/admin"

type AdminHandler struct {
	adminService admin.AdminAdvertisementService
}

func NewAdminHandler(adminService admin.AdminAdvertisementService) *AdminHandler {
	return &AdminHandler{
		adminService: adminService,
	}
}

type RejectionRequest struct {
	Reason string `json:"reason" binding:"required"`
}
