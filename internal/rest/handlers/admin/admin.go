package admin

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *AdminHandler) GetAllAds(c *gin.Context) {
	ads, err := h.adminService.GetAllAds(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get ads: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ads": ads})
}

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
	c.JSON(http.StatusOK, gin.H{"message": "ad deleted"})
}

func (h *AdminHandler) DeleteImage(c *gin.Context) {
	// TODO: implement
}

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
	c.JSON(http.StatusOK, gin.H{"message": "ad approved"})
}

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

	c.JSON(http.StatusOK, gin.H{"message": "ad rejected"})
}
