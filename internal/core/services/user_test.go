package services_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/wansanjou/backend-exercise-user-api/internal/core/domains"
	"github.com/wansanjou/backend-exercise-user-api/internal/core/ports/mocks"
	"github.com/wansanjou/backend-exercise-user-api/internal/core/services"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestUserService_CreateUser(t *testing.T) {
	mockRepo := mocks.NewUserRepository(t)
	userService := services.NewUserService(mockRepo)

	ctx := context.Background()

	inputUser := domains.User{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "hashedpassword123",
	}

	expectedUser := &domains.User{
		ID:        primitive.NewObjectID(),
		Name:      "John Doe",
		Email:     "john@example.com",
		Password:  "hashedpassword123",
		CreatedAt: time.Now(),
	}

	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("domains.User")).
		Return(expectedUser, nil)

	result, err := userService.CreateUser(ctx, inputUser)

	assert.NoError(t, err)
	assert.Equal(t, expectedUser, result)
	assert.Equal(t, expectedUser.Name, result.Name)
	assert.Equal(t, expectedUser.Email, result.Email)

	mockRepo.AssertExpectations(t)
}

func TestUserService_CreateUser_Error(t *testing.T) {
	mockRepo := mocks.NewUserRepository(t)
	userService := services.NewUserService(mockRepo)

	ctx := context.Background()

	inputUser := domains.User{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "hashedpassword123",
	}

	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("domains.User")).
		Return(nil, assert.AnError)

	result, err := userService.CreateUser(ctx, inputUser)

	assert.Error(t, err)
	assert.Nil(t, result)

	mockRepo.AssertExpectations(t)
}

func TestUserService_GetUserByID(t *testing.T) {
	mockRepo := mocks.NewUserRepository(t)
	userService := services.NewUserService(mockRepo)

	ctx := context.Background()

	userID := primitive.NewObjectID()
	expectedUser := &domains.User{
		ID:        primitive.NewObjectID(),
		Name:      "John Doe",
		Email:     "john@example.com",
		Password:  "hashedpassword123",
		CreatedAt: time.Now(),
	}

	mockRepo.On("GetByID", mock.Anything, userID).
		Return(expectedUser, nil)

	result, err := userService.GetUserByID(ctx, userID.Hex())

	assert.NoError(t, err)
	assert.Equal(t, expectedUser, result)
	assert.Equal(t, expectedUser.Name, result.Name)
	assert.Equal(t, expectedUser.Email, result.Email)
	mockRepo.AssertExpectations(t)

}

func TestUserService_GetUserByID_NotFound(t *testing.T) {
	mockRepo := mocks.NewUserRepository(t)
	userService := services.NewUserService(mockRepo)

	ctx := context.Background()

	userID := primitive.NewObjectID()

	mockRepo.On("GetByID", mock.Anything, userID).
		Return(nil, errors.New("user not found"))

	result, err := userService.GetUserByID(ctx, userID.Hex())

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, errors.New("user not found"), err)

	mockRepo.AssertExpectations(t)
}

func TestUserService_GetUsers(t *testing.T) {
	mockRepo := mocks.NewUserRepository(t)
	userService := services.NewUserService(mockRepo)

	ctx := context.Background()

	findParams := domains.FindAllUsers{
		Page:  1,
		Limit: 10,
	}

	expectedUsers := []domains.User{
		{
			ID:        primitive.NewObjectID(),
			Name:      "John Doe",
			Email:     "john@example.com",
			Password:  "hashedpassword123",
			CreatedAt: time.Now(),
		},
		{
			ID:        primitive.NewObjectID(),
			Name:      "John Doe2",
			Email:     "john_2@example.com",
			Password:  "hashedpassword456",
			CreatedAt: time.Now(),
		},
	}

	mockRepo.On("GetUsers", mock.Anything, findParams).
		Return(expectedUsers, nil)

	result, err := userService.GetUsers(ctx, findParams)

	assert.NoError(t, err)
	assert.Equal(t, expectedUsers, result)
	assert.Len(t, result, 2)
	assert.Equal(t, "John Doe", result[0].Name)
	assert.Equal(t, "john@example.com", result[0].Email)

	mockRepo.AssertExpectations(t)
}

func TestUserService_GetUsers_Error(t *testing.T) {
	mockRepo := mocks.NewUserRepository(t)
	userService := services.NewUserService(mockRepo)

	ctx := context.Background()

	findParams := domains.FindAllUsers{
		Page:  1,
		Limit: 10,
	}

	mockRepo.On("GetUsers", mock.Anything, findParams).
		Return(nil, assert.AnError)

	result, err := userService.GetUsers(ctx, findParams)

	assert.Error(t, err)
	assert.Nil(t, result)

	mockRepo.AssertExpectations(t)
}

func TestTransfer_Success(t *testing.T) {
	mockRepo := mocks.NewUserRepository(t)
	userService := services.NewUserService(mockRepo)

	ctx := context.Background()
	fromID := primitive.NewObjectID()
	toID := primitive.NewObjectID()
	amount := 50.0

	mockRepo.On("TransferWithTransaction", ctx, fromID, toID, amount).Return(nil)

	err := userService.TransferBalance(ctx, fromID.Hex(), toID.Hex(), amount)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestTransfer_Error(t *testing.T) {
	mockRepo := mocks.NewUserRepository(t)
	userService := services.NewUserService(mockRepo)

	ctx := context.Background()
	fromID := primitive.NewObjectID()
	toID := primitive.NewObjectID()
	amount := 50.0

	expectedErr := fmt.Errorf("transaction failed")

	mockRepo.On("TransferWithTransaction", ctx, fromID, toID, amount).Return(expectedErr)

	err := userService.TransferBalance(ctx, fromID.Hex(), toID.Hex(), amount)

	assert.Error(t, err)
	assert.Equal(t, expectedErr.Error(), err.Error())

	mockRepo.AssertExpectations(t)
}
