package auth

import (
	"api/internal/users"
	"api/pkg/config"
	"api/pkg/errs"
	"api/pkg/logger"
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"time"
)

// Service implements [IService] and stores methods for authentication.
type Service struct {
	users users.IService
}

// Claims stores data about authorization JWT-token.
type Claims struct {
	jwt.RegisteredClaims
	UserID uint `json:"user_id"`
}

// Login finds a user by given email in db, compares passwords and creates JWT-token.
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

// Signup registrates a new user.
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
	cfg := config.Get()
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(cfg.JwtTokenExp)),
		},
		UserID: userID,
	}

	token := jwt.NewWithClaims(cfg.JwtSigningMethod, claims)

	tokenString, err := token.SignedString([]byte(cfg.JwtTokenSecretKey))

	if err != nil {
		logger.Error(ctx, "failed to create JWT token", err)
		return "", err
	}

	return tokenString, nil
}

func validateJWTToken(ctx context.Context, token string) (Claims, error) {
	claims := Claims{}
	cfg := config.Get()
	data, err := jwt.ParseWithClaims(token, &claims,
		func(t *jwt.Token) (interface{}, error) {
			return []byte(cfg.JwtTokenSecretKey), nil
		},
	)

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return claims, errs.ErrTokenExpired
		}
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return claims, errs.ErrInvalidToken
		}
		return claims, err
	}

	if !data.Valid {
		return claims, errs.ErrInvalidToken
	}

	if data.Method.Alg() != cfg.JwtSigningMethod.Alg() {
		logger.Warn(ctx, "JWT Token method mismatch", "expected_method", cfg.JwtSigningMethod, "got_method", data.Method, "data", data, "token", token)
		return claims, errs.ErrInvalidToken
	}
	return claims, nil
}

func (s *Service) hashPassword(ctx context.Context, password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		logger.Error(ctx, "failed to hash password", err)
		return "", err
	}
	return string(hashed), nil
}

func (s *Service) passwordMatches(hash string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// NewService creates a new Auth Service.
func NewService(users users.IService) *Service {
	return &Service{
		users,
	}
}
