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

type UserRepository interface {
	GetAll() ([]domain.User, error)
	GetById(id int64) (domain.User, error)
	CreateUser(user domain.User) error
	GetByCredentials(email string, hashedPassword string) (domain.User, error)
	Replace(id int64, user domain.User) error
	Update(id int64, userInp domain.UserInput) error
	Delete(id int64) error
}

type PasswordHasher interface {
	Hash(password string) (string, error)
}

type Service struct {
	repository UserRepository
	hasher     PasswordHasher
	hmacSecret []byte
}

func NewService(repository UserRepository, hasher PasswordHasher, secret []byte) *Service {
	return &Service{
		repository: repository,
		hasher:     hasher,
		hmacSecret: secret,
	}
}

func (s *Service) GetAll() ([]domain.User, error) {
	users, err := s.repository.GetAll()
	return users, err
}

func (s *Service) GetById(id int64) (domain.User, error) {
	user, err := s.repository.GetById(id)
	return user, err
}

func (s *Service) SignUp(input domain.UserInput) error {
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

func (s *Service) SignIn(signInInput domain.SignInInput) (string, error) {
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

func (s *Service) ParseToken(ctx context.Context, token string) (int64, error) {
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

func (s *Service) Replace(id int64, user domain.User) error {
	err := s.repository.Replace(id, user)
	return err
}

func (s *Service) Update(id int64, userInp domain.UserInput) error {
	err := s.repository.Update(id, userInp)
	return err
}

func (s *Service) Delete(id int64) error {
	err := s.repository.Delete(id)
	return err
}
