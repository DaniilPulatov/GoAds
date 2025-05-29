//nolint:all // testpackage
package user

import (
	"ads-service/internal/domain/entities"
	"ads-service/internal/usecase/user"
	"bytes"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestUserHandler_CreateDraft(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockService := new(user.MockUserService)
		handler := NewUserHandler(mockService)
		defer mockService.AssertExpectations(t)

		mockService.On("CreateDraft", mock.Anything, "123", mock.Anything).
			Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", "123")

		body := `{"title":"Test Ad","description":"Test desc","category_id":1}`
		req := httptest.NewRequest(http.MethodPost, "/ads", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		c.Request = req
		handler.CreateDraft(c)

		assert.Equal(t, 201, w.Code)
		assert.Equal(t, `{"message":"Ad draft created successfully"}`, w.Body.String())

	})

	t.Run("invalid body", func(t *testing.T) {
		mockService := new(user.MockUserService)
		handler := NewUserHandler(mockService)
		defer mockService.AssertExpectations(t)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", "123")

		body := `{"title":""}`
		req := httptest.NewRequest(http.MethodPost, "/ads", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		c.Request = req

		mockService.On("CreateDraft", mock.Anything, "123", mock.Anything).
			Return(errors.New("validation error"))

		handler.CreateDraft(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "failed to create ad")
	})

	t.Run("service error", func(t *testing.T) {
		mockService := new(user.MockUserService)
		handler := NewUserHandler(mockService)
		defer mockService.AssertExpectations(t)

		mockService.On("CreateDraft", mock.Anything, "123", mock.Anything).
			Return(assert.AnError)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", "123")

		body := `{"title":"Test Ad","description":"Test desc","category_id":1}`
		req := httptest.NewRequest(http.MethodPost, "/ads", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		c.Request = req

		handler.CreateDraft(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "failed to create ad")
	})

	t.Run("missing user_id", func(t *testing.T) {
		mockService := new(user.MockUserService)
		handler := NewUserHandler(mockService)
		defer mockService.AssertExpectations(t)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		body := `{"title":"Test Ad","description":"Test desc","category_id":1}`
		req := httptest.NewRequest(http.MethodPost, "/ads", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		c.Request = req

		handler.CreateDraft(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "unauthorized")
	})
}

func TestUserHandler_GetMyAds(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockService := new(user.MockUserService)
		handler := NewUserHandler(mockService)
		defer mockService.AssertExpectations(t)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		req := httptest.NewRequest("GET", "/ads", nil)
		c.Request = req
		c.Set("user_id", "123")

		expectedAds := []entities.Ad{
			{ID: 1, Title: "Ad 1"},
			{ID: 2, Title: "Ad 2"},
		}

		mockService.On("GetMyAds", mock.Anything, "123").
			Return(expectedAds, nil)

		handler.GetMyAds(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Ad 1")
		assert.Contains(t, w.Body.String(), "Ad 2")
	})

	t.Run("service error", func(t *testing.T) {
		mockService := new(user.MockUserService)
		handler := NewUserHandler(mockService)
		defer mockService.AssertExpectations(t)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		req := httptest.NewRequest("GET", "/ads", nil)
		c.Request = req
		c.Set("user_id", "123")

		mockService.On("GetMyAds", mock.Anything, "123").
			Return(nil, assert.AnError)

		handler.GetMyAds(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "failed to get user ads")
	})
}

func TestUserHandler_UpdateMyAd(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockService := new(user.MockUserService)
		handler := NewUserHandler(mockService)
		defer mockService.AssertExpectations(t)

		mockService.On("UpdateMyAd", mock.Anything, "123", mock.Anything).
			Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", "123")
		body := `{"title":"Updated Ad","description":"desc","category_id":1}`
		req := httptest.NewRequest(http.MethodPut, "/ads/1", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		c.Request = req
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		handler.UpdateMyAd(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Ad updated successfully")
	})

	t.Run("invalid ad id", func(t *testing.T) {
		mockService := new(user.MockUserService)
		handler := NewUserHandler(mockService)
		defer mockService.AssertExpectations(t)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", "123")
		body := `{"title":"Updated Ad"}`
		req := httptest.NewRequest(http.MethodPut, "/ads/abc", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		c.Request = req
		c.Params = gin.Params{{Key: "id", Value: "abc"}}

		handler.UpdateMyAd(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "invalid ad ID")
	})

	t.Run("invalid body", func(t *testing.T) {
		mockService := new(user.MockUserService)
		handler := NewUserHandler(mockService)
		defer mockService.AssertExpectations(t)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", "123")
		body := `invalid-json`
		req := httptest.NewRequest(http.MethodPut, "/ads/1", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		c.Request = req
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		handler.UpdateMyAd(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "invalid request body")
	})

	t.Run("service error", func(t *testing.T) {
		mockService := new(user.MockUserService)
		handler := NewUserHandler(mockService)
		defer mockService.AssertExpectations(t)

		mockService.On("UpdateMyAd", mock.Anything, "123", mock.Anything).
			Return(assert.AnError)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", "123")
		body := `{"title":"Updated Ad","description":"desc","category_id":1}`
		req := httptest.NewRequest(http.MethodPut, "/ads/1", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		c.Request = req
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		handler.UpdateMyAd(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "failed to update ad")
	})
}

func TestUserHandler_DeleteMyAd(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockService := new(user.MockUserService)
		handler := NewUserHandler(mockService)
		defer mockService.AssertExpectations(t)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		req := httptest.NewRequest(http.MethodDelete, "/ads/1", nil)
		c.Request = req
		c.Set("user_id", "123")
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		mockService.On("DeleteMyAd", mock.Anything, "123", 1).
			Return(nil)

		handler.DeleteMyAd(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Ad deleted successfully")
	})

	t.Run("invalid ad id", func(t *testing.T) {
		mockService := new(user.MockUserService)
		handler := NewUserHandler(mockService)
		defer mockService.AssertExpectations(t)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", "123")
		c.Params = gin.Params{{Key: "id", Value: "abc"}}

		handler.DeleteMyAd(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "invalid ad ID")
	})

	t.Run("service error", func(t *testing.T) {
		mockService := new(user.MockUserService)
		handler := NewUserHandler(mockService)
		defer mockService.AssertExpectations(t)

		mockService.On("DeleteMyAd", mock.Anything, "123", 1).
			Return(assert.AnError)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		req := httptest.NewRequest(http.MethodDelete, "/ads/1", nil)
		c.Request = req
		c.Set("user_id", "123")
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		handler.DeleteMyAd(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "failed to delete ad")
	})
}

func TestUserHandler_SubmitForModeration(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockService := new(user.MockUserService)
		handler := NewUserHandler(mockService)
		defer mockService.AssertExpectations(t)

		mockService.On("SubmitForModeration", mock.Anything, "123", 1).
			Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/ads/1/submit", nil)
		c.Set("user_id", "123")
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		handler.SubmitForModeration(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Ad submitted for moderation")
	})

	t.Run("invalid ad id", func(t *testing.T) {
		mockService := new(user.MockUserService)
		handler := NewUserHandler(mockService)
		defer mockService.AssertExpectations(t)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/ads/1/submit", nil)
		c.Set("user_id", "123")
		c.Params = gin.Params{{Key: "id", Value: "abc"}}

		handler.SubmitForModeration(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "invalid ad ID")
	})

	t.Run("service error", func(t *testing.T) {
		mockService := new(user.MockUserService)
		handler := NewUserHandler(mockService)
		defer mockService.AssertExpectations(t)

		mockService.On("SubmitForModeration", mock.Anything, "123", 1).
			Return(assert.AnError)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/ads/1/submit", nil)
		c.Set("user_id", "123")
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		handler.SubmitForModeration(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "failed to submit ad for moderation")
	})
}

func TestUserHandler_AddImageToMyAd(t *testing.T) {
	t.Run("no file provided", func(t *testing.T) {
		mockService := new(user.MockUserService)
		handler := &UserHandler{userService: mockService}

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		err := writer.Close()
		if err != nil {
			return
		}

		req := httptest.NewRequest(http.MethodPost, "/ads/1/image", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "id", Value: "1"}}
		c.Set("user_id", "user-1")

		handler.AddImageToMyAd(c)
		assert.Equal(t, 400, w.Code)
		assert.Contains(t, w.Body.String(), "failed to get file from form")
	})

	t.Run("invalid ad id", func(t *testing.T) {
		mockService := new(user.MockUserService)
		handler := &UserHandler{userService: mockService}

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, _ := writer.CreateFormFile("file", "test.png")
		_, _ = io.Copy(part, bytes.NewBufferString("filecontent"))
		err := writer.Close()
		if err != nil {
			return
		}

		req := httptest.NewRequest(http.MethodPost, "/ads/abc/image", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "id", Value: "abc"}}
		c.Set("user_id", "user-1")

		handler.AddImageToMyAd(c)
		assert.Equal(t, 400, w.Code)
		assert.Contains(t, w.Body.String(), "invalid ad ID")
	})

	t.Run("service error", func(t *testing.T) {
		mockService := new(user.MockUserService)
		handler := &UserHandler{userService: mockService}

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, _ := writer.CreateFormFile("file", "test.png")
		_, _ = io.Copy(part, bytes.NewBufferString("filecontent"))
		err := writer.Close()
		if err != nil {
			return
		}

		req := httptest.NewRequest(http.MethodPost, "/ads/1/image", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "id", Value: "1"}}
		c.Set("user_id", "user-1")

		mockService.On("AddImageToMyAd", mock.Anything, "user-1",
			mock.AnythingOfType("*entities.AdFile")).Return(errors.New("service error"))

		handler.AddImageToMyAd(c)
		assert.Equal(t, 500, w.Code)
		assert.Contains(t, w.Body.String(), "failed to add image to ad")
	})
}

func TestUserHandler_DeleteMyAdImage(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockService := new(user.MockUserService)
		handler := &UserHandler{userService: mockService}
		defer mockService.AssertExpectations(t)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", "user-1")
		c.Params = gin.Params{{Key: "id", Value: "1"}, {Key: "fid", Value: "1"}}

		mockService.On("DeleteMyAdImage", mock.Anything, "user-1",
			mock.AnythingOfType("*entities.AdFile")).
			Return(nil)

		req := httptest.NewRequest(http.MethodDelete, "/ads/1/image/1", nil)
		c.Request = req

		handler.DeleteMyAdImage(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Ad image deleted successfully")
	})

	t.Run("invalid ad id", func(t *testing.T) {
		mockService := new(user.MockUserService)
		handler := &UserHandler{userService: mockService}
		defer mockService.AssertExpectations(t)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", "user-1")
		c.Params = gin.Params{{Key: "id", Value: "abc"}, {Key: "image_id", Value: "img-123"}}

		req := httptest.NewRequest(http.MethodDelete, "/ads/abc/image/img-123", nil)
		c.Request = req

		handler.DeleteMyAdImage(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "invalid ad ID")
	})

	t.Run("service error", func(t *testing.T) {
		mockService := new(user.MockUserService)
		handler := &UserHandler{userService: mockService}
		defer mockService.AssertExpectations(t)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("id", "1")
		c.Params = gin.Params{{Key: "id", Value: "1"}, {Key: "fid", Value: "1"}}

		mockService.On("DeleteMyAdImage", mock.Anything, mock.Anything,
			mock.AnythingOfType("*entities.AdFile")).
			Return(errors.New("some error"))
		req := httptest.NewRequest(http.MethodDelete, "/ads/1", nil)
		c.Request = req

		handler.DeleteMyAdImage(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "failed to delete ad image")
	})
}

func TestUserHandler_GetMyAdsByFilter(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockService := new(user.MockUserService)
		handler := &UserHandler{userService: mockService}
		defer mockService.AssertExpectations(t)

		expectedAds := []entities.Ad{
			{ID: 1, Title: "Ad 1"},
			{ID: 2, Title: "Ad 2"},
		}
		filter := entities.AdFilter{Status: "active", CategoryID: 2, Limit: 10, Page: 1}

		mockService.On("GetMyAdsByFilter", mock.Anything, "user-1", &filter).
			Return(expectedAds, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		req := httptest.NewRequest(http.MethodGet,
			"/ads/filter?status=active&category=2&limit=10&page=1", nil)
		c.Request = req
		c.Set("user_id", "user-1")

		handler.GetMyAdsByFilter(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Ad 1")
		assert.Contains(t, w.Body.String(), "Ad 2")
	})

	t.Run("service error", func(t *testing.T) {
		mockService := new(user.MockUserService)
		handler := &UserHandler{userService: mockService}
		defer mockService.AssertExpectations(t)

		filter := entities.AdFilter{}
		mockService.On("GetMyAdsByFilter", mock.Anything, "user-1", &filter).
			Return(nil, errors.New("some error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		req := httptest.NewRequest(http.MethodGet, "/ads/filter", nil)
		c.Request = req
		c.Set("user_id", "user-1")

		handler.GetMyAdsByFilter(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "failed to get user ads")
	})
}
