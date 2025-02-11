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

type AuthService struct {
	repository         AuthRepository
	sessionsRepository SessionsRepository
	hasher             PasswordHasher
	hmacSecret         []byte
}

func NewAuthService(repository AuthRepository, sessionsRepository SessionsRepository, hasher PasswordHasher, hmacSecret []byte) *AuthService {
	return &AuthService{
		repository:         repository,
		sessionsRepository: sessionsRepository,
		hasher:             hasher,
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
	fmt.Println("inside sign up")
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
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return s.hmacSecret, nil
	})

	if err != nil {
		return 0, err
	}

	if !t.Valid {
		return 0, errors.New("invalid token")
	}

	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid claims")
	}

	subject, ok := claims["sub"].(string)
	if !ok {
		return 0, errors.New("Invalid subject")
	}

	id, err := strconv.Atoi(subject)
	if err != nil {
		return 0, errors.New("invalid subject")
	}

	return int64(id), nil
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
