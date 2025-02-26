package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/dankru/Commissions_simple/internal/domain"
	"github.com/golang-jwt/jwt"
	"math/rand"
	"strconv"
	"time"
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
	ParseToken() (int64, error)
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

func (s *AuthService) SignIn(signInInput domain.SignInInput) (string, string, error) {
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

	return s.generateToken(user.ID)
}

func (s *AuthService) generateToken(userId int64) (string, string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   strconv.Itoa(int(userId)),
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Hour * 15).Unix(),
	})

	accessToken, err := t.SignedString(s.hmacSecret)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := newRefreshToken()
	if err != nil {
		return "", "", err
	}

	if err := s.sessionsRepository.Create(domain.RefreshSession{
		UserID:    userId,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 30),
	}); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *AuthService) ParseToken(ctx context.Context, token string) (int64, error) {
	id, err := s.grpcClient.ParseToken()
	if err != nil {
		return 0, fmt.Errorf(err.Error())
	}

	return id, nil
}

func (s *AuthService) RefreshTokens(refreshToken string) (string, string, error) {
	session, err := s.sessionsRepository.Get(refreshToken)
	if err != nil {
		return "", "", err
	}

	if session.ExpiresAt.Unix() < time.Now().Unix() {
		return "", "", errors.New("refresh token has expired")
	}

	return s.generateToken(session.UserID)
}

func newRefreshToken() (string, error) {
	b := make([]byte, 32)

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	if _, err := r.Read(b); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}
