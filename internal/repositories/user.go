package repositories

import (
	"context"
	"time"

	"github.com/wansanjou/backend-exercise-user-api/internal/core/domains"
	"github.com/wansanjou/backend-exercise-user-api/internal/core/ports"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type userRepository struct {
	mc  *mongo.Client
	db  string
	col string
}

func NewUserRepository(mc *mongo.Client, db string) ports.UserRepository {
	col := "users"
	_, err := mc.Database(db).Collection(col).Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		panic(err)
	}
	return &userRepository{mc, db, col}
}

func (u *userRepository) Create(ctx context.Context, data domains.User) (*domains.User, error) {
	return u.insertOne(ctx, data)
}

func (u *userRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*domains.User, error) {
	return u.findOne(ctx, bson.D{{Key: "_id", Value: id}})
}

func (u *userRepository) GetUsers(ctx context.Context, data domains.FindAllUsers) ([]domains.User, error) {
	filter := bson.D{}
	if data.Name != "" {
		filter = append(filter, bson.E{
			Key: "name", Value: bson.D{{Key: "$regex", Value: data.Name}, {Key: "$options", Value: "i"}},
		})
	}
	if data.Email != "" {
		filter = append(filter, bson.E{
			Key: "email", Value: bson.D{{Key: "$regex", Value: data.Email}, {Key: "$options", Value: "i"}},
		})
	}

	if data.Page < 1 {
		data.Page = 1
	}
	if data.Limit < 1 {
		data.Limit = 10
	}
	skip := int64((data.Page - 1) * data.Limit)
	limit := int64(data.Limit)

	opts := options.Find().SetSkip(skip).SetLimit(limit)

	return u.find(ctx, filter, opts)
}

func (u *userRepository) Count(ctx context.Context) (int64, error) {
	col := u.mc.Database(u.db).Collection(u.col)
	count, err := col.CountDocuments(ctx, bson.D{})
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (u *userRepository) FindByEmail(ctx context.Context, email string) (*domains.User, error) {
	return u.findOne(ctx, bson.D{{Key: "email", Value: email}})
}

func (u *userRepository) TransferWithTransaction(ctx context.Context, fromID, toID primitive.ObjectID, amount float64) error {
	filterFrom := bson.D{
		{Key: "_id", Value: fromID},
		{Key: "balance", Value: bson.D{{Key: "$gte", Value: amount}}},
	}
	updateFrom := bson.D{{Key: "$inc", Value: bson.D{{Key: "balance", Value: -amount}}}}

	fromUser, err := u.findOneAndUpdate(ctx, filterFrom, updateFrom)
	if err != nil {
		return err
	}
	if fromUser == nil {
		return err
	}

	filterTo := bson.D{{Key: "_id", Value: toID}}
	updateTo := bson.D{{Key: "$inc", Value: bson.D{{Key: "balance", Value: amount}}}}

	toUser, err := u.findOneAndUpdate(ctx, filterTo, updateTo)
	if err != nil {
		rollbackFilter := bson.D{{Key: "_id", Value: fromID}}
		rollbackUpdate := bson.D{{Key: "$inc", Value: bson.D{{Key: "balance", Value: amount}}}}
		u.findOneAndUpdate(ctx, rollbackFilter, rollbackUpdate)

		return err
	}
	if toUser == nil {
		rollbackFilter := bson.D{{Key: "_id", Value: fromID}}
		rollbackUpdate := bson.D{{Key: "$inc", Value: bson.D{{Key: "balance", Value: amount}}}}
		u.findOneAndUpdate(ctx, rollbackFilter, rollbackUpdate)

		return err
	}

	return nil
}

func (u *userRepository) find(ctx context.Context, filter bson.D, opts *options.FindOptions) ([]domains.User, error) {
	var out []domains.User
	col := u.mc.Database(u.db).Collection(u.col)

	cursor, err := col.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (u *userRepository) findOne(ctx context.Context, filter bson.D) (*domains.User, error) {
	out := domains.User{}
	col := u.mc.Database(u.db).Collection(u.col)
	if err := col.FindOne(ctx, filter).Decode(&out); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &out, nil
}

func (u *userRepository) insertOne(ctx context.Context, in domains.User) (*domains.User, error) {
	in.CreatedAt = time.Now().UTC()
	col := u.mc.Database(u.db).Collection(u.col)
	result, err := col.InsertOne(ctx, in)
	if err != nil {
		return nil, err
	}
	oid, _ := result.InsertedID.(primitive.ObjectID)
	in.ID = oid
	return &in, nil
}

func (u *userRepository) findOneAndUpdate(ctx context.Context, filter bson.D, update bson.D) (*domains.User, error) {
	col := u.mc.Database(u.db).Collection(u.col)
	var out = domains.User{}
	opts := options.FindOneAndUpdate()
	opts.SetReturnDocument(options.After)
	err := col.FindOneAndUpdate(ctx, filter, update, opts).Decode(&out)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &out, nil
}
