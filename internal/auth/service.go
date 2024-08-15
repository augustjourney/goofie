package auth

import (
	"api/internal/users"
	"api/pkg/config"
	"api/pkg/errs"
	"api/pkg/logger"
	"context"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type Service struct {
	users users.IService
}

type Claims struct {
	jwt.RegisteredClaims
	UserID uint `json:"user_id"`
}

func (s *Service) Login(ctx context.Context, payload LoginDTO) (LoginResult, error) {
	var result LoginResult

	user, err := s.users.GetOneByEmail(ctx, payload.Email)
	if err != nil {
		return result, err
	}

	if !s.passwordMatches(user.Password, payload.Password) {
		return result, errs.ErrWrongCredentials
	}

	token, err := s.createJWTToken(ctx, user.ID)
	if err != nil {
		return result, err
	}

	result.Token = token

	return result, nil
}

func (s *Service) Signup(ctx context.Context, payload SignupDTO) (SignupResult, error) {
	var result SignupResult

	user := users.User{
		Email:     payload.Email,
		Password:  payload.Password,
		Username:  payload.Username,
		FirstName: payload.FirstName,
		LastName:  &payload.LastName,
	}

	hashedPassword, err := s.hashPassword(ctx, payload.Password)
	if err != nil {
		return result, err
	}

	user.Password = hashedPassword

	user, alreadyExists, err := s.users.Create(ctx, user)
	if err != nil {
		return result, err
	}

	result.AlreadyExists = alreadyExists

	return result, nil
}

func (s *Service) createJWTToken(ctx context.Context, userID uint) (string, error) {
	cfg := config.GetConfig()
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(cfg.JwtTokenExp)),
		},
		UserID: userID,
	}

	token := jwt.NewWithClaims(cfg.JwtSigningMethod, claims)

	tokenString, err := token.SignedString([]byte(cfg.JwtTokenSecretKey))

	if err != nil {
		logger.Error(logger.Record{
			Error:   err,
			Context: ctx,
			Message: "[AuthService.createJWTToken] Failed to create JWT token",
		})
		return "", err
	}

	return tokenString, nil
}

func (s *Service) hashPassword(ctx context.Context, password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		logger.Error(logger.Record{
			Context: ctx,
			Error:   err,
			Message: "[AuthService.HashPassword] Failed to hash password",
		})
		return "", err
	}
	return string(hashed), nil
}

func (s *Service) passwordMatches(hash string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func NewService(users users.IService) *Service {
	return &Service{
		users,
	}
}
