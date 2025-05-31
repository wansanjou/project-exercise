package services_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/wansanjou/backend-exercise-user-api/internal/core/domains"
	"github.com/wansanjou/backend-exercise-user-api/internal/core/ports/mocks"
	"github.com/wansanjou/backend-exercise-user-api/internal/core/services"
	"github.com/wansanjou/backend-exercise-user-api/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestAuthService_Login_Success(t *testing.T) {
	mockRepo := mocks.NewUserRepository(t)
	authService := services.NewAuthService(mockRepo)

	ctx := context.Background()

	rawPassword := "hashedpassword456"
	hashedPassword, _ := utils.HashPassword(rawPassword)

	mockUser := &domains.User{
		ID:        primitive.NewObjectID(),
		Name:      "John Doe",
		Email:     "john@example.com",
		Password:  hashedPassword,
		CreatedAt: time.Now(),
	}

	loginInput := domains.LoginRequest{
		Email:    "john@example.com",
		Password: rawPassword,
	}

	mockRepo.On("FindByEmail", mock.Anything, loginInput.Email).Return(mockUser, nil)

	resp, err := authService.Login(ctx, loginInput)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.Token)

	mockRepo.AssertExpectations(t)
}

func TestAuthService_Login_InvalidPassword(t *testing.T) {
	mockRepo := mocks.NewUserRepository(t)
	authService := services.NewAuthService(mockRepo)

	ctx := context.Background()

	rawPassword := "testpassword123"
	hashedPassword, _ := utils.HashPassword("hashedpassword456")

	mockUser := &domains.User{
		ID:        primitive.NewObjectID(),
		Name:      "John Doe",
		Email:     "john@example.com",
		Password:  hashedPassword,
		CreatedAt: time.Now(),
	}

	loginInput := domains.LoginRequest{
		Email:    "john@example.com",
		Password: rawPassword,
	}

	mockRepo.On("FindByEmail", mock.Anything, loginInput.Email).Return(mockUser, nil)

	resp, err := authService.Login(ctx, loginInput)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, "invalid email or password", err.Error())

	mockRepo.AssertExpectations(t)
}
