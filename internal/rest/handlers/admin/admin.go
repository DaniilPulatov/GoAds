package admin

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetAllAds godoc
// @Summary Get all ads
// @Description Get all ads in the system (admin only)
// @Tags admin
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /admin/ads [get]
// @Security BearerAuth
func (h *AdminHandler) GetAllAds(c *gin.Context) {
	ads, err := h.adminService.GetAllAds(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get ads: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ads": ads})
}

// GetStatistics godoc
// @Summary Get statistics
// @Description Get aggregated system statistics (admin only)
// @Tags admin
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /admin/statistics [get]
// @Security BearerAuth
func (h *AdminHandler) GetStatistics(c *gin.Context) {
	stats, err := h.adminService.GetStatistics(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get statistics: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"statistics": stats})
}

func (h *AdminHandler) DeleteAd(c *gin.Context) {
	adID, err := strconv.Atoi(c.Param("id"))
	if err != nil || adID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid ad id"})
		return
	}

	if err := h.adminService.DeleteAd(c.Request.Context(), adID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to delete ad: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Ad deleted"})
}

func (h *AdminHandler) DeleteImage(c *gin.Context) {
	// TODO: implement
}

// Approve godoc
// @Summary Approve ad
// @Description Approve ad by ID (admin only)
// @Tags admin
// @Param id path int true "Ad ID"
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/ads/{id}/approve [post]
// @Security BearerAuth
func (h *AdminHandler) Approve(c *gin.Context) {
	adID, err := strconv.Atoi(c.Param("id"))
	if err != nil || adID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ad id"})
		return
	}

	if err := h.adminService.Approve(c.Request.Context(), adID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to approve ad: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Ad approved"})
}

// Reject godoc
// @Summary Reject ad
// @Description Reject ad by ID with a reason (admin only)
// @Tags admin
// @Param id path int true "Ad ID"
// @Accept json
// @Produce json
// @Param rejection body RejectionRequest true "Rejection reason"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/ads/{id}/reject [post]
// @Security BearerAuth
func (h *AdminHandler) Reject(c *gin.Context) {
	var req RejectionRequest
	adID, err := strconv.Atoi(c.Param("id"))
	if err != nil || adID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ad id"})
		return
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "rejection reason required"})
		return
	}

	if err := h.adminService.Reject(c.Request.Context(), adID, req.Reason); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to reject ad: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ad rejected"})
}
