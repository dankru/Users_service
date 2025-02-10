package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/dankru/Commissions_simple/internal/domain"
	"github.com/golang-jwt/jwt"
	"strconv"
	"time"
)

type AuthRepository interface {
	CreateUser(user domain.User) error
	GetByCredentials(email string, hashedPassword string) (domain.User, error)
}

type AuthService struct {
	repository AuthRepository
	hasher     PasswordHasher
	hmacSecret []byte
}

func NewAuthService(repository AuthRepository, hasher PasswordHasher, hmacSecret []byte) *AuthService {
	return &AuthService{
		repository: repository,
		hasher:     hasher,
		hmacSecret: hmacSecret,
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

func (s *AuthService) SignIn(signInInput domain.SignInInput) (string, error) {
	password, err := s.hasher.Hash(signInInput.Password)
	if err != nil {
		return "", err
	}
	user, err := s.repository.GetByCredentials(signInInput.Email, password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", domain.ErrUserNotFound
		}
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   strconv.Itoa(int(user.ID)),
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Hour * 15).Unix(),
	})
	return token.SignedString(s.hmacSecret)
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
