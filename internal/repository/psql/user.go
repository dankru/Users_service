package psql

import (
	"database/sql"
	"fmt"
	"github.com/dankru/Commissions_simple/internal/domain"
	"strings"
	"time"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (repo *Repository) GetAll() ([]domain.User, error) {
	rows, err := repo.db.Query("select * from users")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	users := make([]domain.User, 0)
	for rows.Next() {
		u := domain.User{}
		err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.RegisteredAt)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (repo *Repository) GetById(id int64) (domain.User, error) {
	var u domain.User
	err := repo.db.QueryRow("select * from users WHERE id = $1", id).
		Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.RegisteredAt)
	return u, err
}

func (repo *Repository) CreateUser(user domain.User) error {
	_, err := repo.db.Exec("insert into users (name, email, password, registered_at) values ($1, $2, $3, $4)",
		user.Name, user.Email, user.Password, time.Now())
	return err
}

func (repo *Repository) GetByCredentials(email string, hashedPassword string) (domain.User, error) {
	var user domain.User
	err := repo.db.QueryRow("SELECT id, name, email, password, registered_at FROM users WHERE email=$1 AND password=$2",
		email, hashedPassword).
		Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.RegisteredAt)

	return user, err
}

func (repo *Repository) Replace(id int64, user domain.User) error {
	_, err := repo.db.Exec("update users set name = $1, email = $2, password = $3 WHERE id = $4",
		user.Name, user.Email, user.Password, id)
	return err
}

func (repo *Repository) Update(id int64, userInp domain.UserInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if userInp.Name != nil {
		setValues = append(setValues, fmt.Sprintf("name = $%d", argId))
		args = append(args, userInp.Name)
		argId++
	}

	if userInp.Email != nil {
		setValues = append(setValues, fmt.Sprintf("email = $%d", argId))
		args = append(args, userInp.Email)
		argId++
	}

	if userInp.Password != nil {
		setValues = append(setValues, fmt.Sprintf("password = $%d", argId))
		args = append(args, userInp.Password)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf("update users set %s where id = $%d", setQuery, argId)

	args = append(args, id)

	_, err := repo.db.Exec(query, args...)

	return err
}

func (repo *Repository) Delete(id int64) error {
	_, err := repo.db.Exec("delete from users where id = $1", id)
	return err
}
