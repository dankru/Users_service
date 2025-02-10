package service

import (
	"github.com/dankru/Commissions_simple/internal/domain"
)

type UserRepository interface {
	GetAll() ([]domain.User, error)
	GetById(id int64) (domain.User, error)
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

func NewService(repository UserRepository) *Service {
	return &Service{
		repository: repository,
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
