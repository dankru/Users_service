package rest

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/dankru/Commissions_simple/internal/domain"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"strconv"
)

type CtxValue int

const (
	ctxUserId CtxValue = iota
)

type AuthService interface {
	SignUp(user domain.UserInput) error
	SignIn(ctx context.Context, signInInput domain.SignInInput) (string, string, error)
	ParseToken(ctx context.Context, token string) (int64, error)
	RefreshTokens(ctx context.Context, refreshToken string) (string, string, error)
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
	r.StrictSlash(true)
	r.Use(loggingMiddleware)
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

func decodeJsonBody[T domain.Input](r *http.Request) (T, error) {
	var dst T

	reqBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return dst, err
	}

	if err = json.Unmarshal(reqBytes, &dst); err != nil {
		return dst, err
	}
	return dst, nil
}
