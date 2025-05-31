package services

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
	"github.com/wansanjou/backend-exercise-user-api/internal/core/domains"
	"github.com/wansanjou/backend-exercise-user-api/internal/core/ports"
	"github.com/wansanjou/backend-exercise-user-api/utils"
)

type authService struct {
	userrepo ports.UserRepository
}

func NewAuthService(userrepo ports.UserRepository) ports.AuthService {
	return &authService{
		userrepo: userrepo,
	}
}

func (s *authService) Login(ctx context.Context, in domains.LoginRequest) (*domains.LoginResponse, error) {
	if in.Email == "" {
		return nil, errors.New("email is required")
	}
	if in.Password == "" {
		return nil, errors.New("password is required")
	}

	user, err := s.userrepo.FindByEmail(ctx, in.Email)
	if err != nil || user == nil {
		return nil, errors.New("invalid email or password")
	}

	if err := utils.VerifyPassword(in.Password, user.Password); err != nil {
		return nil, errors.New("invalid email or password")
	}

	claims := domains.JWTClaims{
		ID:    user.ID.Hex(),
		Email: user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secretKey := viper.GetString("jwt.secretKey")
	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return nil, err
	}

	return &domains.LoginResponse{
		Token: signedToken,
	}, nil
}
