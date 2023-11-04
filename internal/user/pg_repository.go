package user

import (
	"context"
	"github.com/blazee5/testhub-backend/internal/models"
)

type Repository interface {
	GetUserById(ctx context.Context, userId int) (models.User, error)
}
