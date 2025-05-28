package auth

import (
	"ads-service/internal/usecase/auth"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAuthHandler_Register(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockService := new(auth.MockAuthService)
		handler := NewAuthHandler(mockService)
		defer mockService.AssertExpectations(t)

		body := `{"username":"testuser","password":"testpass"}`
		req := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		mockService.On("Register", mock.Anything, mock.Anything).Return(nil)

		handler.Register(c)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Equal(t, `{"message":"user registered successfully"}`, w.Body.String())
	})

	t.Run("invalid username", func(t *testing.T) {
		mockService := new(auth.MockAuthService)
		handler := NewAuthHandler(mockService)
		defer mockService.AssertExpectations(t)

		body := `{"username":testuser,"password":123}`
		req := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		handler.Register(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "invalid request body")
	})

	t.Run("invalid password", func(t *testing.T) {
		mockService := new(auth.MockAuthService)
		handler := NewAuthHandler(mockService)
		defer mockService.AssertExpectations(t)

		body := `{"username":"testuser","password":123}`
		req := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		handler.Register(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "invalid request body")
	})

	t.Run("fail service", func(t *testing.T) {
		mockService := new(auth.MockAuthService)
		handler := NewAuthHandler(mockService)
		defer mockService.AssertExpectations(t)

		body := `{"username":"testuser","password":"testpass"}`
		req := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		mockService.On("Register", mock.Anything, mock.Anything).Return(assert.AnError)

		handler.Register(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "failed to register user")
	})
}

func TestAuthHandler_Login_Success(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockService := new(auth.MockAuthService)
		handler := NewAuthHandler(mockService)
		defer mockService.AssertExpectations(t)

		body := `{"phone":"1234567890","password":"testpass"}`
		req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		mockService.On("Login", mock.Anything, "1234567890", "testpass").Return("access-token", "refresh-token", nil)

		handler.Login(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "access-token")
		assert.Contains(t, w.Body.String(), "refresh-token")
	})

	t.Run("invalid phone", func(t *testing.T) {
		mockService := new(auth.MockAuthService)
		handler := NewAuthHandler(mockService)
		defer mockService.AssertExpectations(t)

		body := `{"phone":1234567890,"password":"testpass"}` // некорректный тип поля
		req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		handler.Login(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "invalid request body")
	})

	t.Run("invalid password", func(t *testing.T) {
		mockService := new(auth.MockAuthService)
		handler := NewAuthHandler(mockService)
		defer mockService.AssertExpectations(t)

		body := `{"phone":"1234567890","password":123}` // некорректный тип поля
		req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		handler.Login(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "invalid request body")
	})

	t.Run("invalid body", func(t *testing.T) {
		mockService := new(auth.MockAuthService)
		handler := NewAuthHandler(mockService)
		defer mockService.AssertExpectations(t)

		body := `{"phone":123}` // некорректный тип поля
		req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		handler.Login(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "invalid request body")
	})

	t.Run("fail service", func(t *testing.T) {
		mockService := new(auth.MockAuthService)
		handler := NewAuthHandler(mockService)
		defer mockService.AssertExpectations(t)

		body := `{"phone":"1234567890","password":"testpass"}`
		req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		mockService.On("Login", mock.Anything, "1234567890", "testpass").Return("", "", assert.AnError)

		handler.Login(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "failed to login")
	})

}
