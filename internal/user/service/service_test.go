package service

import (
	"context"
	"github.com/blazee5/quizmaster-backend/internal/models"
	mock_user "github.com/blazee5/quizmaster-backend/internal/user/mock"
	"github.com/blazee5/quizmaster-backend/lib/logger"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestService_GetById(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	user := models.UserInfo{
		User: models.ShortUser{
			Id:       1,
			Username: "username",
			Email:    "email@gmail.com",
			Avatar:   "",
		},
		Quizzes: []models.Quiz{},
		Results: []models.UserResult{},
	}

	ctx := context.Background()

	log := logger.NewLogger()
	mockUserRepo := mock_user.NewMockRepository(ctrl)
	mockUserRedisRepo := mock_user.NewMockRedisRepository(ctrl)
	userService := NewService(log, mockUserRepo, mockUserRedisRepo)

	mockUserRedisRepo.EXPECT().GetByIdCtx(ctx, "1").Return(&user, nil)

	selectedUser, err := userService.GetById(ctx, 1)

	require.NoError(t, err)
	require.NotNil(t, selectedUser)
	require.Nil(t, err)
	require.Equal(t, selectedUser, user)
}
