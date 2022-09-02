package service

import (
	"auth/internal/domain/models"
	"auth/internal/ports"
	"auth/internal/utils"
	"context"
)

type UserService struct {
	repo ports.UserRepo
}

func NewUserService(repo ports.UserRepo) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (us *UserService) GetAll(ctx context.Context) ([]*models.User, error) {
	ctx, span := utils.StartSpan(ctx)
	defer span.End()

	users, err := us.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (us *UserService) Create(ctx context.Context, user *models.User) error {
	ctx, span := utils.StartSpan(ctx)
	defer span.End()

	hash, err := utils.HashAndSalt(user.Password)
	if err != nil {
		return err
	}
	user.Password = hash
	return us.repo.Insert(ctx, user)
}

func (us *UserService) Update(ctx context.Context, user *models.User) error {
	ctx, span := utils.StartSpan(ctx)
	defer span.End()

	return us.repo.Update(ctx, user)
}

func (us *UserService) UpdatePassword(ctx context.Context, user *models.User) error {
	ctx, span := utils.StartSpan(ctx)
	defer span.End()

	hash, err := utils.HashAndSalt(user.Password)
	if err != nil {
		return err
	}
	user.Password = hash
	return us.repo.UpdatePassword(ctx, user)
}
