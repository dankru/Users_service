package rest

import (
	"context"
	"errors"
	"github.com/dankru/Commissions_simple/internal/domain"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type CtxValue int

const (
	ctxUserId CtxValue = iota
)

type AuthService interface {
	SignUp(user domain.UserInput) error
	SignIn(signInInput domain.SignInInput) (string, error)
	ParseToken(ctx context.Context, token string) (int64, error)
}

type UserService interface {
	GetAll() ([]domain.User, error)
	GetById(id int64) (domain.User, error)
	Replace(id int64, user domain.User) error
	Update(id int64, userInp domain.UserInput) error
	Delete(id int64) error
}

type Handler struct {
	authService AuthService
	userService UserService
}

func NewHandler(authService AuthService, userService UserService) *Handler {
	return &Handler{
		authService: authService,
		userService: userService,
	}
}

func (h *Handler) InitRouter() *mux.Router {
	r := mux.NewRouter()
	h.initAuthRoutes(r)
	h.initUserRoutes(r)
	return r
}

func getIdFromRequest(r *http.Request) (int64, error) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		return 0, err
	}

	if id == 0 {
		return 0, errors.New("id can't be 0")
	}

	return id, nil
}
