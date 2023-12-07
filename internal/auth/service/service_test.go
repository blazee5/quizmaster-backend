package service

import (
	"context"
	"github.com/blazee5/quizmaster-backend/internal/auth/mock"
	"github.com/blazee5/quizmaster-backend/internal/domain"
	"github.com/blazee5/quizmaster-backend/internal/models"
	"github.com/blazee5/quizmaster-backend/lib/auth"
	"github.com/blazee5/quizmaster-backend/lib/logger"
	"github.com/blazee5/quizmaster-backend/lib/tracer"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestSignUp(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	user := domain.SignUpRequest{
		Username: "username",
		Email:    "email@gmail.com",
		Password: "123456",
	}

	ctx := context.Background()

	log := logger.NewLogger()
	mockAuthRepo := mock_auth.NewMockRepository(ctrl)
	authService := NewService(log, mockAuthRepo, tracer.InitTracer("main", "test"))

	mockAuthRepo.EXPECT().CreateUser(gomock.Any(), gomock.Eq(domain.SignUpRequest{
		Username: "username",
		Email:    "email@gmail.com",
		Password: auth.GenerateHashPassword(user.Password),
	})).Return(0, nil)

	createdUser, err := authService.SignUp(ctx, user)
	require.NoError(t, err)
	require.NotNil(t, createdUser)
	require.Nil(t, err)
}

func TestSignIn(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	user := domain.SignInRequest{
		Email:    "email@gmail.com",
		Password: "123456",
	}

	ctx := context.Background()

	log := logger.NewLogger()

	mockAuthRepo := mock_auth.NewMockRepository(ctrl)
	authService := NewService(log, mockAuthRepo, tracer.InitTracer("main", "test"))

	mockAuthRepo.EXPECT().ValidateUser(gomock.Any(), gomock.Eq(domain.SignInRequest{
		Email:    "email@gmail.com",
		Password: auth.GenerateHashPassword(user.Password),
	})).Return(models.User{Email: user.Email, Password: user.Password}, nil)

	token, err := authService.GenerateToken(ctx, user)
	require.NoError(t, err)
	require.NotNil(t, token)
	require.Nil(t, err)
}
