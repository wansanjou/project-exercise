package domains

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name"`
	Email     string             `bson:"email"`
	Password  string             `bson:"password"`
	Balance   float64            `bson:"balance"`
	CreatedAt time.Time          `bson:"created_at"`
}

type CreateUser struct {
	Name      string
	Email     string
	Password  string
	CreatedAt time.Time
}

type FindAllUsers struct {
	Name  string
	Email string
	Page  int
	Limit int
}

type TransferRequest struct {
	FromUserID string  `json:"fromUserId"`
	ToUserID   string  `json:"toUserId"`
	Amount     float64 `json:"amount"`
}
