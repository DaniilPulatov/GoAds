package auth

import (
	"ads-service/internal/domain/entities"
	"ads-service/internal/repository/auth"
	"ads-service/internal/repository/user"
	customLogger "ads-service/pkg/logger"
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestMockAuthService_IsAdmin(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockUserRepo := &user.MockUserRepo{}
		mockAuthRepo := &auth.MockAuthRepository{}
		defer mockUserRepo.AssertExpectations(t)
		defer mockAuthRepo.AssertExpectations(t)

		service := NewAuthService(mockUserRepo, mockAuthRepo, customLogger.Logger{})

		mockUserRepo.On("GetUserByID", mock.Anything, "1").
			Return(&entities.User{Role: entities.RoleAdmin}, nil)

		isAdmin, err := service.IsAdmin(context.Background(), "1")
		assert.NoError(t, err)
		assert.True(t, isAdmin)
	})

	t.Run("not admin", func(t *testing.T) {
		mockUserRepo := &user.MockUserRepo{}
		mockAuthRepo := &auth.MockAuthRepository{}
		defer mockUserRepo.AssertExpectations(t)
		defer mockAuthRepo.AssertExpectations(t)

		service := NewAuthService(mockUserRepo, mockAuthRepo, customLogger.Logger{})

		mockUserRepo.On("GetUserByID", mock.Anything, "2").
			Return(&entities.User{Role: entities.RoleUser}, nil)

		isAdmin, err := service.IsAdmin(context.Background(), "2")
		assert.NoError(t, err)
		assert.False(t, isAdmin)
	})

	t.Run("user not found", func(t *testing.T) {
		mockUserRepo := &user.MockUserRepo{}
		mockAuthRepo := &auth.MockAuthRepository{}
		defer mockUserRepo.AssertExpectations(t)
		defer mockAuthRepo.AssertExpectations(t)

		service := NewAuthService(mockUserRepo, mockAuthRepo, customLogger.Logger{})

		mockUserRepo.On("GetUserByID", mock.Anything, "3").Return(nil, nil)

		isAdmin, err := service.IsAdmin(context.Background(), "3")
		assert.Error(t, err)
		assert.False(t, isAdmin)
	})

	t.Run("repo error", func(t *testing.T) {
		mockUserRepo := &user.MockUserRepo{}
		mockAuthRepo := &auth.MockAuthRepository{}
		defer mockUserRepo.AssertExpectations(t)
		defer mockAuthRepo.AssertExpectations(t)

		service := NewAuthService(mockUserRepo, mockAuthRepo, customLogger.Logger{})

		mockUserRepo.On("GetUserByID", mock.Anything, "4").Return(nil, assert.AnError)

		isAdmin, err := service.IsAdmin(context.Background(), "4")
		assert.Error(t, err)
		assert.False(t, isAdmin)
	})
}

func TestMockAuthService_Register(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockUserRepo := &user.MockUserRepo{}
		mockAuthRepo := &auth.MockAuthRepository{}
		defer mockUserRepo.AssertExpectations(t)
		defer mockAuthRepo.AssertExpectations(t)

		service := NewAuthService(mockUserRepo, mockAuthRepo, customLogger.Logger{})
		userEntity := &entities.User{Phone: "+998917773355", Password: "ValidPass123!"}

		mockUserRepo.On("IsExists", mock.Anything, userEntity.Phone).
			Return(false, nil)
		mockUserRepo.On("CreateUser", mock.Anything, mock.AnythingOfType("*entities.User")).
			Return("1", nil)

		err := service.Register(context.Background(), userEntity)
		assert.NoError(t, err)
	})

	t.Run("empty pass and phone", func(t *testing.T) {
		mockUserRepo := &user.MockUserRepo{}
		mockAuthRepo := &auth.MockAuthRepository{}
		service := NewAuthService(mockUserRepo, mockAuthRepo, customLogger.Logger{})
		userEntity := &entities.User{Phone: "", Password: ""}

		err := service.Register(context.Background(), userEntity)
		assert.Error(t, err)
	})

	t.Run("incorrect pass or number", func(t *testing.T) {
		mockUserRepo := &user.MockUserRepo{}
		mockAuthRepo := &auth.MockAuthRepository{}
		service := NewAuthService(mockUserRepo, mockAuthRepo, customLogger.Logger{})
		userEntity := &entities.User{Phone: "123", Password: "123"}

		err := service.Register(context.Background(), userEntity)
		assert.Error(t, err)
	})

	t.Run("user already exist", func(t *testing.T) {
		mockUserRepo := &user.MockUserRepo{}
		mockAuthRepo := &auth.MockAuthRepository{}
		service := NewAuthService(mockUserRepo, mockAuthRepo, customLogger.Logger{})
		userEntity := &entities.User{Phone: "+998917773355", Password: "ValidPass123!"}

		mockUserRepo.On("IsExists", mock.Anything, userEntity.Phone).Return(true, nil)

		err := service.Register(context.Background(), userEntity)
		assert.Error(t, err)
	})

	t.Run("err checking user", func(t *testing.T) {
		mockUserRepo := &user.MockUserRepo{}
		mockAuthRepo := &auth.MockAuthRepository{}
		service := NewAuthService(mockUserRepo, mockAuthRepo, customLogger.Logger{})
		userEntity := &entities.User{Phone: "+998917773355", Password: "ValidPass123!"}

		mockUserRepo.On("IsExists", mock.Anything, userEntity.Phone).
			Return(false, assert.AnError)

		err := service.Register(context.Background(), userEntity)
		assert.Error(t, err)
	})

	t.Run("err creation user", func(t *testing.T) {
		mockUserRepo := &user.MockUserRepo{}
		mockAuthRepo := &auth.MockAuthRepository{}
		service := NewAuthService(mockUserRepo, mockAuthRepo, customLogger.Logger{})
		userEntity := &entities.User{Phone: "+998917773355", Password: "ValidPass123!"}

		mockUserRepo.On("IsExists", mock.Anything, userEntity.Phone).Return(false, nil)
		mockUserRepo.On("CreateUser", mock.Anything, mock.AnythingOfType("*entities.User")).Return("", assert.AnError)

		err := service.Register(context.Background(), userEntity)
		assert.Error(t, err)
	})
}

func TestMockAuthService_Login(t *testing.T) {
	t.Setenv("REFRESH_TOKEN_LIFETIME", "10")
	t.Setenv("ACCESS_TOKEN_LIFETIME", "5")

	t.Run("success login", func(t *testing.T) {
		mockUserRepo := &user.MockUserRepo{}
		mockAuthRepo := &auth.MockAuthRepository{}
		defer mockUserRepo.AssertExpectations(t)
		defer mockAuthRepo.AssertExpectations(t)

		service := NewAuthService(mockUserRepo, mockAuthRepo, customLogger.Logger{})
		password := "ValidPass123!"
		hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		userEntity := &entities.User{ID: "1", Phone: "+79999999999", PasswordHash: string(hashed)}

		mockUserRepo.On("GetByPhone", mock.Anything, userEntity.Phone).Return(userEntity, nil)
		mockAuthRepo.On("Create", mock.Anything, mock.AnythingOfType("entities.Token")).
			Return(nil)

		rToken, accessToken, err := service.Login(context.Background(), userEntity.Phone, password)
		assert.NoError(t, err)
		assert.NotEmpty(t, rToken)
		assert.NotEmpty(t, accessToken)
	})

	t.Run("empty phone and password", func(t *testing.T) {
		mockUserRepo := &user.MockUserRepo{}
		mockAuthRepo := &auth.MockAuthRepository{}
		service := NewAuthService(mockUserRepo, mockAuthRepo, customLogger.Logger{})

		rToken, accessToken, err := service.Login(context.Background(), "", "")
		assert.Error(t, err)
		assert.Empty(t, rToken)
		assert.Empty(t, accessToken)
	})

	t.Run("user not found", func(t *testing.T) {
		mockUserRepo := &user.MockUserRepo{}
		mockAuthRepo := &auth.MockAuthRepository{}
		service := NewAuthService(mockUserRepo, mockAuthRepo, customLogger.Logger{})

		mockUserRepo.On("GetByPhone", mock.Anything, "notfound").Return(nil, nil)

		rToken, accessToken, err := service.Login(context.Background(), "notfound", "pass")
		assert.Error(t, err)
		assert.Empty(t, rToken)
		assert.Empty(t, accessToken)
	})

	t.Run("error getting user", func(t *testing.T) {
		mockUserRepo := &user.MockUserRepo{}
		mockAuthRepo := &auth.MockAuthRepository{}
		service := NewAuthService(mockUserRepo, mockAuthRepo, customLogger.Logger{})

		mockUserRepo.On("GetByPhone", mock.Anything, "err").Return(nil, assert.AnError)

		rToken, accessToken, err := service.Login(context.Background(), "err", "pass")
		assert.Error(t, err)
		assert.Empty(t, rToken)
		assert.Empty(t, accessToken)
	})

	t.Run("wrong password", func(t *testing.T) {
		mockUserRepo := &user.MockUserRepo{}
		mockAuthRepo := &auth.MockAuthRepository{}
		service := NewAuthService(mockUserRepo, mockAuthRepo, customLogger.Logger{})
		hashed, _ := bcrypt.GenerateFromPassword([]byte("rightpass"), bcrypt.DefaultCost)
		userEntity := &entities.User{ID: "1", Phone: "+79999999999", PasswordHash: string(hashed)}

		mockUserRepo.On("GetByPhone", mock.Anything, userEntity.Phone).Return(userEntity, nil)

		rToken, accessToken, err := service.Login(context.Background(), userEntity.Phone, "wrongpass")
		assert.Error(t, err)
		assert.Empty(t, rToken)
		assert.Empty(t, accessToken)
	})

	t.Run("error token creation", func(t *testing.T) {
		mockUserRepo := &user.MockUserRepo{}
		mockAuthRepo := &auth.MockAuthRepository{}
		service := NewAuthService(mockUserRepo, mockAuthRepo, customLogger.Logger{})
		password := "ValidPass123!"
		hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		userEntity := &entities.User{ID: "1", Phone: "+79999999999", PasswordHash: string(hashed)}

		mockUserRepo.On("GetByPhone", mock.Anything, userEntity.Phone).Return(userEntity, nil)
		mockAuthRepo.On("Create", mock.Anything, mock.AnythingOfType("entities.Token")).Return(assert.AnError)

		rToken, accessToken, err := service.Login(context.Background(), userEntity.Phone, password)
		assert.Error(t, err)
		assert.Empty(t, rToken)
		assert.Empty(t, accessToken)
	})
}

func TestMockAuthService_Refresh(t *testing.T) {
	t.Setenv("REFRESH_TOKEN_LIFETIME", "10")
	t.Setenv("ACCESS_TOKEN_LIFETIME", "5")
	t.Setenv("JWT_SECRET_KEY", "testsecret")

	t.Run("ошибка конвертации времени", func(t *testing.T) {
		t.Setenv("REFRESH_TOKEN_LIFETIME", "notint")
		mockUserRepo := &user.MockUserRepo{}
		mockAuthRepo := &auth.MockAuthRepository{}
		service := NewAuthService(mockUserRepo, mockAuthRepo, customLogger.Logger{})

		access, refresh, err := service.Refresh(context.Background(), "sometoken")
		assert.Error(t, err)
		assert.Empty(t, access)
		assert.Empty(t, refresh)
	})
}
