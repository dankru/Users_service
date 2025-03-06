package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/dankru/Commissions_simple/internal/domain"
)

type AuthRepository interface {
	CreateUser(user domain.User) error
	GetByCredentials(email string, hashedPassword string) (domain.User, error)
}

type SessionsRepository interface {
	Create(token domain.RefreshSession) error
	Get(token string) (domain.RefreshSession, error)
}

type GrpcClient interface {
	ParseToken(ctx context.Context, token string) (int64, error)
	GenerateToken(ctx context.Context, userId int64) (string, string, error)
	RefreshToken(ctx context.Context, token string) (string, string, error)
}

type AuthService struct {
	repository         AuthRepository
	sessionsRepository SessionsRepository
	hasher             PasswordHasher
	grpcClient         GrpcClient
	hmacSecret         []byte
}

func NewAuthService(repository AuthRepository, sessionsRepository SessionsRepository, hasher PasswordHasher, grpcClient GrpcClient, hmacSecret []byte) *AuthService {
	return &AuthService{
		repository:         repository,
		sessionsRepository: sessionsRepository,
		hasher:             hasher,
		grpcClient:         grpcClient,
		hmacSecret:         hmacSecret,
	}
}

func (s *AuthService) SignUp(input domain.UserInput) error {

	password, err := s.hasher.Hash(*input.Password)
	if err != nil {
		return err
	}
	user := domain.User{
		Name:     *input.Name,
		Email:    *input.Email,
		Password: password,
	}
	err = s.repository.CreateUser(user)
	return err
}

func (s *AuthService) SignIn(ctx context.Context, signInInput domain.SignInInput) (string, string, error) {
	password, err := s.hasher.Hash(signInInput.Password)
	if err != nil {
		return "", "", err
	}
	user, err := s.repository.GetByCredentials(signInInput.Email, password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", "", domain.ErrUserNotFound
		}
		return "", "", err
	}

	return s.GenerateToken(ctx, user.ID)
}

func (s *AuthService) GenerateToken(ctx context.Context, userId int64) (string, string, error) {
	accessToken, refreshToken, err := s.grpcClient.GenerateToken(ctx, userId)
	if err != nil {
		return "", "", fmt.Errorf(err.Error())
	}

	return accessToken, refreshToken, nil
}

func (s *AuthService) ParseToken(ctx context.Context, token string) (int64, error) {
	id, err := s.grpcClient.ParseToken(ctx, token)
	if err != nil {
		return 0, fmt.Errorf(err.Error())
	}

	return id, nil
}

func (s *AuthService) RefreshTokens(ctx context.Context, refreshToken string) (string, string, error) {
	accessToken, refreshToken, err := s.grpcClient.RefreshToken(ctx, refreshToken)
	if err != nil {
		return "", "", fmt.Errorf(err.Error())
	}

	return accessToken, refreshToken, err
}
