package grpc

import (
	"auth/internal/domain/models"
	"auth/internal/ports"
	"auth_grpc"
	"context"
	"errors"
)

type AuthApi struct {
	authService ports.AuthService
	auth_grpc.UnimplementedAuthServiceServer
}

func NewAuthAPI(authS ports.AuthService) *AuthApi {
	return &AuthApi{authService: authS}
}

func (a *AuthApi) Validate(ctx context.Context, req *auth_grpc.ValidateTokenRequest) (*auth_grpc.ValidateTokenResponse, error) {
	user, err := a.authService.ValidateTokens(ctx, &models.TokenPair{
		AccessToken:  req.AccessToken,
		RefreshToken: req.RefreshToken,
	})

	if err != nil {
		return &auth_grpc.ValidateTokenResponse{
			AccessToken:  "",
			RefreshToken: "",
			Status:       auth_grpc.Statuses_invalid,
		}, err
	}

	if errors.Is(err, models.ErrTokenExpired) {
		tokens, err := a.authService.Authorize(ctx, user)
		if err != nil {
			return &auth_grpc.ValidateTokenResponse{
				AccessToken:  "",
				RefreshToken: "",
				Status:       auth_grpc.Statuses_expired,
			}, err
		}
		return &auth_grpc.ValidateTokenResponse{
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
			Status:       auth_grpc.Statuses_expired,
		}, nil
	}

	return &auth_grpc.ValidateTokenResponse{
		AccessToken:  req.GetAccessToken(),
		RefreshToken: req.GetRefreshToken(),
		Status:       auth_grpc.Statuses_valid,
	}, nil
}
