package user

import (
	"ads-service/internal/domain/entities"
	"ads-service/internal/usecase/user"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService user.UserAdvertisementService
}

func NewUserHandler(userService user.UserAdvertisementService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) AddImageToMyAd(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"error": "failed to get file from form: " + err.Error()})
		return
	}

	adID := c.Param("id")
	intID, err := strconv.Atoi(adID)
	if err != nil || intID <= 0 {
		c.JSON(400, gin.H{"error": "invalid ad ID: " + err.Error()})
		return
	}

	adFile := &entities.AdFile{
		FileName: file.Filename,
		AdID:     intID,
	}
	err = h.userService.AddImageToMyAd(c.Request.Context(), c.GetString("user_id"), adFile)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to add image to ad: " + err.Error()})
		return
	}
	if err := c.SaveUploadedFile(file, adFile.URL); err != nil {
		c.JSON(500, gin.H{"error": "failed to save uploaded file: " + err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "image added to ad successfully", "file": adFile})
}

func (h *UserHandler) DeleteMyAdImage(c *gin.Context) {
	adID := c.Param("id")
	fID := c.Param("fid")
	intID, err := strconv.Atoi(adID)
	if err != nil || intID <= 0 {
		c.JSON(400, gin.H{"error": "invalid ad ID: " + err.Error()})
		return
	}
	intfID, err := strconv.Atoi(fID)
	if err != nil || intfID <= 0 {
		c.JSON(400, gin.H{"error": "invalid ad file ID: " + err.Error()})
		return
	}

	file := &entities.AdFile{
		ID:  intfID,
		AdID: intID,
	}

	if err := h.userService.DeleteMyAdImage(c.Request.Context(), c.GetString("user_id"), file); err != nil {
		c.JSON(500, gin.H{"error": "failed to delete ad image: " + err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "ad image deleted successfully"})
}

func (h *UserHandler) GetImagesToMyAd(c *gin.Context) {
	adID := c.Param("id")
	intID, err := strconv.Atoi(adID)
	if err != nil || intID <= 0 {
		c.JSON(400, gin.H{"error": "invalid ad ID: " + err.Error()})
		return
	}

	files, err := h.userService.GetImagesToMyAd(c.Request.Context(), c.GetString("user_id"), intID)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to get ad images: " + err.Error()})
		return
	}
	c.JSON(200, gin.H{"files": files})
}