package ports

import (
	"context"

	"github.com/wansanjou/backend-exercise-user-api/internal/core/domains"
)

type UserService interface {
	CreateUser(ctx context.Context, data domains.User) (*domains.User, error)
	GetUserByID(ctx context.Context, id string) (*domains.User, error)
	GetUsers(ctx context.Context, data domains.FindAllUsers) ([]domains.User, error)
	TransferBalance(ctx context.Context, fromID, toID string, amount float64) error
	CountUsers(ctx context.Context) (int64, error)
}

type AuthService interface {
	Login(ctx context.Context, in domains.LoginRequest) (*domains.LoginResponse, error)
}
