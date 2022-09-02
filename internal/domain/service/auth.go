package service

import (
	"auth/internal/config"
	"auth/internal/domain/models"
	"auth/internal/ports"
	"auth/internal/utils"
	"context"
	"errors"
	"time"
)

type AuthService struct {
	repo ports.UserRepo
	cfg  *config.Configs
}

func NewAuthService(repository ports.UserRepo, configs *config.Configs) *AuthService {
	return &AuthService{
		repo: repository,
		cfg:  configs,
	}
}

func (s *AuthService) Authorize(ctx context.Context, user *models.User) (*models.TokenDetails, error) {
	ctx, span := utils.StartSpan(ctx)
	defer span.End()

	userData, err := s.repo.GetByName(ctx, user.Login)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			return nil, models.ErrUserNotFound
		}
		return nil, err
	}

	if err := utils.ComparePasswordAndHash(userData.Password, user.Password); err != nil {
		return nil, models.ErrInvalidPassword
	}

	tokens, err := s.generateTokens(user)
	if err != nil {
		return nil, err
	}

	return tokens, nil
}

func (s *AuthService) ValidateTokens(ctx context.Context, pair *models.TokenPair) (*models.User, error) {
	ctx, span := utils.StartSpan(ctx)
	defer span.End()

	login, err := utils.ParseToken(pair.AccessToken)
	if err != nil {
		return nil, err
	}
	if login == "" {
		return nil, models.ErrInvalidToken
	}
	_, err = utils.ParseToken(pair.RefreshToken)
	if err != nil {
		return nil, err
	}
	return &models.User{Login: login}, nil

}

func (s *AuthService) generateTokens(user *models.User) (*models.TokenDetails, error) {
	accesToken, err := utils.GenerateToken(user.Login, s.cfg.Jwt.AccessTokenExpTime)
	if err != nil {
		return nil, err
	}
	refreshToken, err := utils.GenerateToken(user.Login, s.cfg.Jwt.RefreshTokenExpTime)
	if err != nil {
		return nil, err
	}

	return &models.TokenDetails{
		TokenPair: models.TokenPair{
			AccessToken:  accesToken,
			RefreshToken: refreshToken,
		},
		AtExpires: time.Now().Add(s.cfg.Jwt.AccessTokenExpTime),
		RtExpires: time.Now().Add(s.cfg.Jwt.RefreshTokenExpTime),
	}, nil
}

func (s *AuthService) SignUp(ctx context.Context, user *models.User) error {
	ctx, span := utils.StartSpan(ctx)
	defer span.End()

	_, err := s.repo.GetByName(ctx, user.Login)
	if err != nil {
		if !errors.Is(err, models.ErrUserNotFound) {
			return err
		}
	} else {
		return models.ErrUserExist
	}

	hashedPass, err := utils.HashAndSalt(user.Password)
	if err != nil {
		return err
	}

	if err := s.repo.Insert(ctx, &models.User{
		Login:        user.Login,
		Password:     hashedPass,
		Email:        user.Email,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Role:         user.Role,
		CreationDate: uint64(time.Now().Unix()),
	}); err != nil {
		return err
	}

	return nil
}
