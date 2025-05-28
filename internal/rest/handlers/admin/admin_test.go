package admin

import (
	"ads-service/internal/domain/entities"
	"ads-service/internal/usecase/admin"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAdminHandler_Approve(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockService := new(admin.MockAdminService)
		handler := NewAdminHandler(mockService)
		defer mockService.AssertExpectations(t)

		mockService.On("Approve", mock.Anything, 1).Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{
			{Key: "id", Value: "1"},
		}

		c.Request = httptest.NewRequest(http.MethodPost, "/admin/ads/1/approve", nil)
		handler.Approve(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, `{"message":"Ad approved"}`, w.Body.String())
	})

	t.Run("invalid id", func(t *testing.T) {
		mockService := new(admin.MockAdminService)
		handler := NewAdminHandler(mockService)
		defer mockService.AssertExpectations(t)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{
			{Key: "id", Value: "abc"},
		}
		c.Request = httptest.NewRequest(http.MethodPost, "/admin/ads/abc/approve", nil)
		handler.Approve(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "invalid ad id")
	})

	t.Run("service error", func(t *testing.T) {
		mockService := new(admin.MockAdminService)
		handler := NewAdminHandler(mockService)
		defer mockService.AssertExpectations(t)

		mockService.On("Approve", mock.Anything, 2).Return(assert.AnError)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{
			{Key: "id", Value: "2"},
		}
		c.Request = httptest.NewRequest(http.MethodPost, "/admin/ads/2/approve", nil)
		handler.Approve(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "failed to approve ad")
	})
}

func TestAdminHandler_Reject(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockService := new(admin.MockAdminService)
		handler := NewAdminHandler(mockService)
		defer mockService.AssertExpectations(t)

		mockService.On("Reject", mock.Anything, 1, "spam").Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{
			{Key: "id", Value: "1"},
		}
		body := `{"rejection_reason":"spam"}`
		c.Request = httptest.NewRequest(http.MethodPost, "/admin/ads/1/reject",
			strings.NewReader(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.Reject(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, `{"message":"Ad rejected"}`, w.Body.String())
	})

	t.Run("invalid id", func(t *testing.T) {
		mockService := new(admin.MockAdminService)
		handler := NewAdminHandler(mockService)
		defer mockService.AssertExpectations(t)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{
			{Key: "id", Value: "abc"},
		}
		body := `{"reason":"spam"}`
		c.Request = httptest.NewRequest(http.MethodPost, "/admin/ads/abc/reject",
			strings.NewReader(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.Reject(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "invalid ad id")
	})

	t.Run("invalid body", func(t *testing.T) {
		mockService := new(admin.MockAdminService)
		handler := NewAdminHandler(mockService)
		defer mockService.AssertExpectations(t)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{
			{Key: "id", Value: "1"},
		}
		body := `{}`
		c.Request = httptest.NewRequest(http.MethodPost, "/admin/ads/1/reject",
			strings.NewReader(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.Reject(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "rejection reason required")
	})

	t.Run("service error", func(t *testing.T) {
		mockService := new(admin.MockAdminService)
		handler := NewAdminHandler(mockService)
		defer mockService.AssertExpectations(t)

		mockService.On("Reject", mock.Anything, 1, "spam").Return(assert.AnError)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{
			{Key: "id", Value: "1"},
		}
		body := `{"rejection_reason":"spam"}`
		c.Request = httptest.NewRequest(http.MethodPost, "/admin/ads/1/reject",
			strings.NewReader(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.Reject(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "failed to reject ad")
	})
}

func TestAdminHandler_GetAllAds(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockService := new(admin.MockAdminService)
		handler := NewAdminHandler(mockService)
		defer mockService.AssertExpectations(t)

		mockService.On("GetAllAds", mock.Anything).
			Return([]entities.Ad{
				{Title: "ad1", Description: "desc1", CategoryID: 1},
				{Title: "ad2", Description: "desc2", CategoryID: 2},
			}, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/admin/ads", nil)

		handler.GetAllAds(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "ad1")
		assert.Contains(t, w.Body.String(), "ad2")
	})

	t.Run("error", func(t *testing.T) {
		mockService := new(admin.MockAdminService)
		handler := NewAdminHandler(mockService)
		defer mockService.AssertExpectations(t)

		mockService.On("GetAllAds", mock.Anything).
			Return([]entities.Ad{}, assert.AnError)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/admin/ads", nil)

		handler.GetAllAds(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "failed to get ads")
	})

}

func TestAdminHandler_GetStatistics(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockService := new(admin.MockAdminService)
		handler := NewAdminHandler(mockService)
		defer mockService.AssertExpectations(t)

		mockService.On("GetStatistics", mock.Anything).
			Return(entities.AdStatistics{Total: 10, Published: 5, Pending: 3, Rejected: 2}, nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/admin/statistics", nil)

		handler.GetStatistics(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"Total":10`)
	})

	t.Run("error", func(t *testing.T) {
		mockService := new(admin.MockAdminService)
		handler := NewAdminHandler(mockService)
		defer mockService.AssertExpectations(t)

		mockService.On("GetStatistics", mock.Anything).
			Return(entities.AdStatistics{}, assert.AnError)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/admin/statistics", nil)

		handler.GetStatistics(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "failed to get statistics")
	})
}

func TestAdminHandler_DeleteAd(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockService := new(admin.MockAdminService)
		handler := NewAdminHandler(mockService)
		defer mockService.AssertExpectations(t)

		mockService.On("DeleteAd", mock.Anything, 1).Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{
			{Key: "id", Value: "1"},
		}
		c.Request = httptest.NewRequest(http.MethodDelete, "/admin/ads/1", nil)

		handler.DeleteAd(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Ad deleted")
	})

	t.Run("invalid id", func(t *testing.T) {
		mockService := new(admin.MockAdminService)
		handler := NewAdminHandler(mockService)
		defer mockService.AssertExpectations(t)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{
			{Key: "id", Value: "abc"},
		}
		c.Request = httptest.NewRequest(http.MethodDelete, "/admin/ads/abc", nil)

		handler.DeleteAd(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "invalid ad id")
	})

	t.Run("service error", func(t *testing.T) {
		mockService := new(admin.MockAdminService)
		handler := NewAdminHandler(mockService)
		defer mockService.AssertExpectations(t)

		mockService.On("DeleteAd", mock.Anything, 2).Return(assert.AnError)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{
			{Key: "id", Value: "2"},
		}
		c.Request = httptest.NewRequest(http.MethodDelete, "/admin/ads/2", nil)

		handler.DeleteAd(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "failed to delete ad")
	})
}
