package user

import (
	"ads-service/internal/domain/entities"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateDraft godoc
// @Summary      Create a new ad draft
// @Description  Allows a user to create an ad draft
// @Tags         user-ads
// @Accept       json
// @Produce      json
// @Param        ad  body      entities.Ad  true  "Ad draft"
// @Success      201  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Security BearerAuth
// @Router       /ads [post]
func (h *UserHandler) CreateDraft(c *gin.Context) {
	var ad entities.Ad
	if err := c.ShouldBindJSON(&ad); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return
	}

	userID := c.GetString("user_id")
	if err := h.userService.CreateDraft(c.Request.Context(), userID, &ad); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create ad: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "ad draft created successfully"})
}

// GetMyAds godoc
// @Summary      Get all ads created by the authenticated user
// @Tags         user-ads
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]string
// @Security BearerAuth
// @Router       /ads [get]
func (h *UserHandler) GetMyAds(c *gin.Context) {
	userID := c.GetString("user_id")
	ads, err := h.userService.GetMyAds(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user ads: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ads": ads})
}

// UpdateMyAd godoc
// @Summary      Update a user's own ad
// @Tags         user-ads
// @Accept       json
// @Produce      json
// @Param        id   path      int          true  "Ad ID"
// @Param        ad   body      entities.Ad  true  "Updated ad data"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Security BearerAuth
// @Router       /ads/{id} [put]
func (h *UserHandler) UpdateMyAd(c *gin.Context) {
	var ad entities.Ad
	adIDStr := c.Param("id")
	adID, err := strconv.Atoi(adIDStr)
	if err != nil || adID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ad ID: " + err.Error()})
		return
	}

	ad.ID = adID

	if err := c.ShouldBindJSON(&ad); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return
	}

	userID := c.GetString("user_id")
	if err := h.userService.UpdateMyAd(c.Request.Context(), userID, &ad); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update ad: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ad updated successfully"})
}

// DeleteMyAd godoc
// @Summary      Delete a user's own ad
// @Tags         user-ads
// @Produce      json
// @Param        id   path      int  true  "Ad ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Security BearerAuth
// @Router       /ads/{id} [delete]
func (h *UserHandler) DeleteMyAd(c *gin.Context) {
	adIDStr := c.Param("id")
	adID, err := strconv.Atoi(adIDStr)
	if err != nil || adID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ad ID: " + err.Error()})
		return
	}

	userID := c.GetString("user_id")
	if err := h.userService.DeleteMyAd(c.Request.Context(), userID, adID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete ad: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ad deleted successfully"})
}

// SubmitForModeration godoc
// @Summary      Submit an ad for moderation
// @Tags         user-ads
// @Produce      json
// @Param        id   path      int  true  "Ad ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Security BearerAuth
// @Router       /ads/{id}/submit [post]
func (h *UserHandler) SubmitForModeration(c *gin.Context) {
	adIDStr := c.Param("id")
	adID, err := strconv.Atoi(adIDStr)
	if err != nil || adID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ad ID: " + err.Error()})
		return
	}

	userID := c.GetString("user_id")
	if err := h.userService.SubmitForModeration(c.Request.Context(), userID, adID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to submit ad for moderation: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ad submitted for moderation"})
}

// AddImageToMyAd godoc
// @Summary Add image to user's ad
// @Description Uploads and attaches an image file to the user's draft ad
// @Tags user-ads
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "Ad ID"
// @Param file formData file true "Image file"
// @Success 200 {object} map[string]interface{} "image added successfully"
// @Failure 400 {object} map[string]string "invalid input"
// @Failure 500 {object} map[string]string "internal error"
// @Security BearerAuth
// @Router /ads/{id}/image [post]
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
	log.Println("ADID IS : ", adID)
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

// DeleteMyAdImage godoc
// @Summary Delete image from user's ad
// @Description Deletes a specific image from the user's ad by ad ID and file ID
// @Tags user-ads
// @Produce json
// @Param id path int true "Ad ID"
// @Param fid path int true "File ID"
// @Success 200 {object} map[string]string "image deleted successfully"
// @Failure 400 {object} map[string]string "invalid input"
// @Failure 500 {object} map[string]string "internal error"
// @Security BearerAuth
// @Router /ads/{id}/image/{fid} [delete]
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
		ID:   intfID,
		AdID: intID,
	}

	if err := h.userService.DeleteMyAdImage(c.Request.Context(), c.GetString("user_id"), file); err != nil {
		c.JSON(500, gin.H{"error": "failed to delete ad image: " + err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "ad image deleted successfully"})
}

// GetImagesToMyAd godoc
// @Summary Get images of user's ad
// @Description Returns all images attached to the user's ad by ID
// @Tags user-ads
// @Produce json
// @Param id path int true "Ad ID"
// @Success 200 {object} map[string]interface{} "list of images"
// @Failure 400 {object} map[string]string "invalid input"
// @Failure 500 {object} map[string]string "internal error"
// @Security BearerAuth
// @Router /ads/{id}/image [get]
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
