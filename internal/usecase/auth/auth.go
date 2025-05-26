package auth

import (
	"ads-service/internal/domain/entities"
	repoerr "ads-service/internal/errs/repoErr"
	usecaserr "ads-service/internal/errs/usecaseErr"
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

func (s *authService) Register(ctx context.Context, user entities.User) error {
	if user.Password == "" || user.Phone == "" {
		log.Println("HERE!!!!")
		return usecaserr.ErrInvalidUserData
	}
	if len(user.Password) < minPswdLenfth || !utils.IsValidPhone(user.Phone) {
		log.Println("invalid phone or password")
		return usecaserr.ErrInvalidUserData
	}
	isExists, err := s.userRepo.IsExists(ctx, user.Phone)
	if err != nil {
		log.Println("Error checking if user exists:", err)
		return usecaserr.ErrCheckUserExists
	}
	if isExists {
		log.Println("User already exists with phone:", user.Phone)
		return usecaserr.ErrUserAlreadyExists
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		log.Println("Error hashing password:", err)
		return usecaserr.ErrInvalidUserData
	}
	user.PasswordHash = string(hash)
	user.Password = "" // Clear the password field after hashing

	err = s.userRepo.Create(ctx, user)
	if err != nil {
		log.Println("Error creating user:", err)
		return repoerr.ErrUserInsertFailed
	}
	log.Println("User registered successfully:", user.Phone)

	return nil
}

func (s *authService) Login(ctx context.Context, phone, password string) (string, string, error) {
	if phone == "" || password == "" {
		return "", "", usecaserr.ErrInvalidUserData
	}

	user, err := s.userRepo.GetByPhone(ctx, phone)
	if err != nil {
		log.Println("Error getting user by phone:", err)
		return "", "", usecaserr.ErrCheckUserExists
	}

	if user == nil {
		log.Println("User not found with phone:", phone)
		return "", "", usecaserr.ErrUserNotFound
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		log.Println("Password mismatch for user:", phone)
		return "", "", usecaserr.ErrInvalidUserData
	}

	RToken, err := utils.GenerateToken(user.ID, refershTokenDuration)
	if err != nil {
		log.Println("Error generating refresh token:", err)
		return "", "", usecaserr.ErrTokenGeneration
	}
	accessToken, err := utils.GenerateToken(user.ID, accessTokenDuration)
	if err != nil {
		log.Println("Error generating access token:", err)
		return "", "", usecaserr.ErrTokenGeneration
	}

	if err := s.authRepo.CreateToken(ctx, entities.RefreshToken{
		UserID:    user.ID,
		Token:     RToken,
		ExpiresAt: time.Now().Local().Add(refershTokenDuration * time.Minute),
	}); err != nil {
		log.Println("Error creating refresh token in repository:", err)
		return "", "", usecaserr.ErrTokenGeneration
	}
	return RToken, accessToken, nil
}

func (s *authService) Refresh(ctx context.Context, refreshToken string) (string, string, error) {
	claims := &utils.CustomClaims{}
	token, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return utils.SecretKey, nil
	})
	if err != nil {
		log.Println("err while parsing claims: ", err)
		return "", "", usecaserr.ErrInvalidToken
	}
	if !token.Valid {
		log.Println("token is not valid")
		return "", "", usecaserr.ErrInvalidToken
	}
	if err := s.authRepo.DeleteToken(ctx, refreshToken); err != nil {
		log.Println("err while deleting refresh token: ", err)
		return "", "", usecaserr.ErrInvalidToken
	}
	newAccessToken, err := utils.GenerateToken(claims.UserID, accessTokenDuration)
	if err != nil {
		log.Println(err)
		return "", "", usecaserr.ErrTokenGeneration
	}
	newRefreshToken, err := utils.GenerateToken(claims.UserID, refershTokenDuration)
	if err != nil {
		log.Println(err)
		return "", "", usecaserr.ErrTokenGeneration
	}
	err = s.authRepo.CreateToken(ctx, entities.RefreshToken{
		UserID:    claims.UserID,
		Token:     newRefreshToken,
		ExpiresAt: time.Now().Local().Add(refershTokenDuration * time.Second),
	})
	if err != nil {
		log.Println(err)
		return "", "", usecaserr.ErrTokenGeneration
	}
	return newAccessToken, newRefreshToken, nil
}
