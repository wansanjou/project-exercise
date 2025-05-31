package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/wansanjou/backend-exercise-user-api/internal/core/domains"
	"github.com/wansanjou/backend-exercise-user-api/internal/core/ports"
	"github.com/wansanjou/backend-exercise-user-api/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type service struct {
	userrepo ports.UserRepository
}

func NewUserService(userrepo ports.UserRepository) ports.UserService {
	return &service{
		userrepo: userrepo,
	}
}

func (s *service) CreateUser(ctx context.Context, data domains.User) (*domains.User, error) {
	if data.Name == "" {
		return nil, errors.New("name is required")
	}
	if data.Email == "" {
		return nil, errors.New("email is required")
	}
	if data.Password == "" {
		return nil, errors.New("password is required")
	}

	data.Balance = 100.0

	hash, err := utils.HashPassword(data.Password)
	if err != nil {
		return nil, err
	}
	data.Password = hash

	return s.userrepo.Create(ctx, data)
}

func (s *service) GetUserByID(ctx context.Context, id string) (*domains.User, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid id")
	}
	return s.userrepo.GetByID(ctx, oid)
}

func (s *service) GetUsers(ctx context.Context, data domains.FindAllUsers) ([]domains.User, error) {
	return s.userrepo.GetUsers(ctx, data)
}

func (s *service) CountUsers(ctx context.Context) (int64, error) {
	return s.userrepo.Count(ctx)
}

func (s *service) TransferBalance(ctx context.Context, fromID, toID string, amount float64) error {
	if fromID == toID {
		return fmt.Errorf("cannot transfer to the same user")
	}

	if amount <= 0.00 {
		return fmt.Errorf("amount must be greater than zero")
	}

	foid, err := primitive.ObjectIDFromHex(fromID)
	if err != nil {
		return fmt.Errorf("invalid from user ID: %v", err)
	}

	toid, err := primitive.ObjectIDFromHex(toID)
	if err != nil {
		return fmt.Errorf("invalid to user ID: %v", err)
	}

	return s.userrepo.TransferWithTransaction(ctx, foid, toid, amount)
}
