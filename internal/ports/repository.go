package ports

import (
	"auth/internal/domain/models"
	"context"
)

type UserRepo interface {
	Get(ctx context.Context, id string) (*models.User, error)
	GetAll(ctx context.Context) ([]*models.User, error)
	GetByName(ctx context.Context, login string) (*models.User, error)
	Insert(ctx context.Context, user *models.User) error
	Update(ctx context.Context, user *models.User) error
	UpdatePassword(ctx context.Context, user *models.User) error
}
