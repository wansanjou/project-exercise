package ports

import (
	"context"

	"github.com/wansanjou/backend-exercise-user-api/internal/core/domains"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRepository interface {
	Create(ctx context.Context, data domains.User) (*domains.User, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (*domains.User, error)
	GetUsers(ctx context.Context, data domains.FindAllUsers) ([]domains.User, error)
	Count(ctx context.Context) (int64, error)
	TransferWithTransaction(ctx context.Context, fromID, toID primitive.ObjectID, amount float64) error

	//Auth
	FindByEmail(ctx context.Context, email string) (*domains.User, error)
}
