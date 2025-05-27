package auth

import (
	"ads-service/internal/domain/entities"
	repoerr "ads-service/internal/errs/repoErr"
	usecaseerr "ads-service/internal/errs/usecaseErr"
	"ads-service/pkg/utils"
	"context"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// TODO: Move to config or env variable

const (
	refershTokenDuration = 10_080  // Refersh Token 7 days duration in minutes
	accessTokenDuration  = 60 * 15 // Access Token15 min duration in minutes
	minPswdLenfth        = 8
) // Token duration in minutes

func (s *userAuthService) Register(ctx context.Context, user *entities.User) error {
	if user.Password == "" || user.Phone == "" {
		log.Println("HERE!!!!")
		return usecaseerr.ErrInvalidUserData
	}
	if len(user.Password) < minPswdLenfth || !utils.IsValidPhone(user.Phone) {
		log.Println("invalid phone or password")
		return usecaseerr.ErrInvalidUserData
	}
	isExists, err := s.userRepo.IsExists(ctx, user.Phone)
	if err != nil {
		log.Println("Error checking if user exists:", err)
		return usecaseerr.ErrCheckUserExists
	}
	if isExists {
		log.Println("User already exists with phone:", user.Phone)
		return usecaseerr.ErrUserAlreadyExists
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		log.Println("Error hashing password:", err)
		return usecaseerr.ErrInvalidUserData
	}
	user.PasswordHash = string(hash)
	user.Password = "" // Clear the password field after hashing

	_, err = s.userRepo.CreateUser(ctx, user)
	if err != nil {
		log.Println("Error creating user:", err)
		return repoerr.ErrUserInsertFailed
	}
	log.Println("User registered successfully:", user.Phone)

	return nil
}

func (s *userAuthService) Login(ctx context.Context, phone, password string) (rToken, accessToken string, err error) {
	if phone == "" || password == "" {
		return "", "", usecaseerr.ErrInvalidUserData
	}

	user, err := s.userRepo.GetByPhone(ctx, phone)
	if err != nil {
		log.Println("Error getting user by phone:", err)
		return "", "", usecaseerr.ErrCheckUserExists
	}

	if user == nil {
		log.Println("User not found with phone:", phone)
		return "", "", usecaseerr.ErrUserNotFound
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		log.Println("Password mismatch for user:", phone)
		return "", "", usecaseerr.ErrInvalidUserData
	}

	rToken, err = utils.GenerateToken(user.ID, refershTokenDuration)
	if err != nil {
		log.Println("Error generating refresh token:", err)
		return "", "", usecaseerr.ErrTokenGeneration
	}
	accessToken, err = utils.GenerateToken(user.ID, accessTokenDuration)
	if err != nil {
		log.Println("Error generating access token:", err)
		return "", "", usecaseerr.ErrTokenGeneration
	}

	if err := s.authRepo.CreateToken(ctx, entities.RefreshToken{
		UserID:    user.ID,
		Token:     rToken,
		ExpiresAt: time.Now().Local().Add(refershTokenDuration * time.Minute),
	}); err != nil {
		log.Println("Error creating refresh token in repository:", err)
		return "", "", usecaseerr.ErrTokenGeneration
	}
	return rToken, accessToken, nil
}

func (s *userAuthService) Refresh(ctx context.Context, refreshToken string) (newAccessToken, newRefreshToken string, err error) {
	claims := &utils.CustomClaims{}
	token, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return utils.SecretKey, nil
	})
	if err != nil {
		log.Println("err while parsing claims: ", err)
		return "", "", usecaseerr.ErrInvalidToken
	}
	if !token.Valid {
		log.Println("token is not valid")
		return "", "", usecaseerr.ErrInvalidToken
	}
	if err := s.authRepo.DeleteToken(ctx, refreshToken); err != nil {
		log.Println("err while deleting refresh token: ", err)
		return "", "", usecaseerr.ErrInvalidToken
	}
	newAccessToken, err = utils.GenerateToken(claims.UserID, accessTokenDuration)
	if err != nil {
		log.Println(err)
		return "", "", usecaseerr.ErrTokenGeneration
	}
	newRefreshToken, err = utils.GenerateToken(claims.UserID, refershTokenDuration)
	if err != nil {
		log.Println(err)
		return "", "", usecaseerr.ErrTokenGeneration
	}
	err = s.authRepo.CreateToken(ctx, entities.RefreshToken{
		UserID:    claims.UserID,
		Token:     newRefreshToken,
		ExpiresAt: time.Now().Local().Add(refershTokenDuration * time.Second),
	})
	if err != nil {
		log.Println(err)
		return "", "", usecaseerr.ErrTokenGeneration
	}
	return newAccessToken, newRefreshToken, nil
}

func (s *userAuthService) IsAdmin(ctx context.Context, userID string) (bool, error) {
	userByID, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return false, usecaseerr.ErrGettingUser
	}
	if userByID == nil {
		return false, usecaseerr.ErrUserNotFound
	}

	return userByID.Role == entities.RoleAdmin, nil
}
