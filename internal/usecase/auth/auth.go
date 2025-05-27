package auth

import (
	"ads-service/internal/domain/entities"
	repoerr "ads-service/internal/errs/repoerr"
	usecaseerr "ads-service/internal/errs/usecaseerr"

	"ads-service/pkg/utils"
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// TODO: Move to config or env variable

func (s *userAuthService) Register(ctx context.Context, user *entities.User) error {
	if user.Password == "" || user.Phone == "" {
		log.Println("phone or password is empty")
		return usecaseerr.ErrInvalidUserData
	}
	if !utils.IsValidPassword(user.Password) || !utils.IsValidPhone(user.Phone) {
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
	intRefresh, err := strconv.Atoi(os.Getenv("REFRESH_TOKEN_LIFETIME"))
	if err != nil {
		log.Println("Error converting refresh token duration to int:", err)
		return "", "", usecaseerr.ErrInvalidTokenDuration
	}

	intAccess, err := strconv.Atoi(os.Getenv("ACCESS_TOKEN_LIFETIME"))
	if err != nil {
		log.Println("Error converting access token duration to int:", err)
		return "", "", usecaseerr.ErrInvalidTokenDuration
	}

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

	rToken, err = utils.GenerateToken(user.ID, intRefresh)
	if err != nil {
		log.Println("Error generating refresh token:", err)
		return "", "", usecaseerr.ErrTokenGeneration
	}
	accessToken, err = utils.GenerateToken(user.ID, intAccess)
	if err != nil {
		log.Println("Error generating access token:", err)
		return "", "", usecaseerr.ErrTokenGeneration
	}

	if err := s.authRepo.Create(ctx, entities.Token{
		UserID:    user.ID,
		Token:     rToken,
		ExpiresAt: time.Now().Local().Add(time.Duration(intRefresh) * time.Minute),
	}); err != nil {
		log.Println("Error creating refresh token in repository:", err)
		return "", "", usecaseerr.ErrTokenGeneration
	}
	return rToken, accessToken, nil
}

func (s *userAuthService) Refresh(
	ctx context.Context,
	refreshToken string,
) (newAccessToken, newRefreshToken string, err error) {
	intRefresh, err := strconv.Atoi(os.Getenv("REFRESH_TOKEN_LIFETIME"))
	if err != nil {
		log.Println("Error converting refresh token duration to int:", err)
		return "", "", usecaseerr.ErrInvalidTokenDuration
	}

	intAccess, err := strconv.Atoi(os.Getenv("ACCESS_TOKEN_LIFETIME"))
	if err != nil {
		log.Println("Error converting access token duration to int:", err)
		return "", "", usecaseerr.ErrInvalidTokenDuration
	}

	claims := &utils.CustomClaims{}
	token, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return os.Getenv("JWT_SECRET_KEY"), nil
	})
	if err != nil {
		log.Println("err while parsing claims: ", err)
		return "", "", usecaseerr.ErrInvalidToken
	}
	if !token.Valid {
		log.Println("token is not valid")
		return "", "", usecaseerr.ErrInvalidToken
	}
	if err := s.authRepo.Delete(ctx, refreshToken); err != nil {
		log.Println("err while deleting refresh token: ", err)
		return "", "", usecaseerr.ErrInvalidToken
	}
	newAccessToken, err = utils.GenerateToken(claims.UserID, intAccess)
	if err != nil {
		log.Println(err)
		return "", "", usecaseerr.ErrTokenGeneration
	}
	newRefreshToken, err = utils.GenerateToken(claims.UserID, intRefresh)
	if err != nil {
		log.Println(err)
		return "", "", usecaseerr.ErrTokenGeneration
	}
	err = s.authRepo.Create(ctx, entities.Token{
		UserID:    claims.UserID,
		Token:     newRefreshToken,
		ExpiresAt: time.Now().Local().Add(time.Duration(intAccess) * time.Minute),
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
