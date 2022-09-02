package ports

import (
	"auth/internal/domain/models"
	"context"
)

type AuthService interface {
	Authorize(ctx context.Context, user *models.User) (*models.TokenDetails, error)
	ValidateTokens(ctx context.Context, tokens *models.TokenPair) (*models.User, error)
}

type UserService interface {
	GetAll(ctx context.Context) ([]*models.User, error)
	Create(ctx context.Context, user *models.User) (err error)
	Update(ctx context.Context, user *models.User) (err error)
	UpdatePassword(ctx context.Context, user *models.User) (err error)
}
