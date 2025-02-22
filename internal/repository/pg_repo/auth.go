package pg_repo

import (
	"database/sql"
	"github.com/dankru/Commissions_simple/internal/domain"
	"time"
)

type AuthRepository struct {
	db *sql.DB
}

type UserRepository interface {
	CreateUser(user domain.User) error
	GetByCredentials(email string, hashedPassword string) (domain.User, error)
}

func NewAuthRepository(db *sql.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

func (repo *AuthRepository) CreateUser(user domain.User) error {
	_, err := repo.db.Exec("insert into users (name, email, password, registered_at) values ($1, $2, $3, $4)",
		user.Name, user.Email, user.Password, time.Now())
	return err
}

func (repo *AuthRepository) GetByCredentials(email string, hashedPassword string) (domain.User, error) {
	var user domain.User
	err := repo.db.QueryRow("SELECT id, name, email, password, registered_at FROM users WHERE email=$1 AND password=$2",
		email, hashedPassword).
		Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.RegisteredAt)

	return user, err
}
