package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"strconv"
)

type UserService interface {
	GetAll() ([]domain.User, error)
	GetById(id int64) (domain.User, error)
	SignUp(user domain.UserInput) error
	Replace(id int64, user domain.User) error
	Update(id int64, userInp domain.UserInput) error
	Delete(id int64) error
}

type Handler struct {
	service UserService
}

func NewHandler(service UserService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) InitRouter() *mux.Router {
	r := mux.NewRouter()
	h.initAuthRoutes(r)
	h.initUserRoutes(r)
	return r
}

func (h *Handler) initAuthRoutes(router *mux.Router) {
	users := router.PathPrefix("/auth").Subrouter()
	{
		users.HandleFunc("/sign-up", h.signUp).Methods(http.MethodPost)
	}
}

func (h *Handler) initUserRoutes(router *mux.Router) {
	users := router.PathPrefix("/users").Subrouter()
	{
		users.HandleFunc("", h.getUsers).Methods(http.MethodGet)
		users.HandleFunc("/{id:[0-9]+}", h.getUserById).Methods(http.MethodGet)
		users.HandleFunc("/{id:[0-9]+}", h.replaceUser).Methods(http.MethodPut)
		users.HandleFunc("/{id:[0-9]+}", h.updateUser).Methods(http.MethodPatch)
		users.HandleFunc("/{id:[0-9]+}", h.deleteUser).Methods(http.MethodDelete)
	}
}

func (h *Handler) getUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.GetAll()
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get users: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(users)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to marshall users: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	w.Write([]byte(resp))
}

func (h *Handler) getUserById(w http.ResponseWriter, r *http.Request) {
	id, err := getIdFromRequest(r)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	user, err := h.service.GetById(id)
	if err != nil {
		http.Error(w, "failed to find user", http.StatusInternalServerError)
		return
	}
	response, err := json.Marshal(user)
	if err != nil {
		http.Error(w, "failed to marshall user", http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(response)
}

func (h *Handler) signUp(w http.ResponseWriter, r *http.Request) {

	reqBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	var user domain.UserInput
	if err = json.Unmarshal(reqBytes, &user); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	if err = user.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = h.service.SignUp(user); err != nil {
		http.Error(w, fmt.Sprintf("failed to create user: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) replaceUser(w http.ResponseWriter, r *http.Request) {
	id, err := getIdFromRequest(r)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	reqBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	var user domain.User
	if err := json.Unmarshal(reqBytes, &user); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	if err := h.service.Replace(id, user); err != nil {
		http.Error(w, fmt.Sprintf("failed to update user: %s", err.Error()), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) updateUser(w http.ResponseWriter, r *http.Request) {
	id, err := getIdFromRequest(r)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	reqBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	var userInp domain.UserInput
	if err := json.Unmarshal(reqBytes, &userInp); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	if err := h.service.Update(id, userInp); err != nil {
		http.Error(w, fmt.Sprintf("failed to update user: %s", err.Error()), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) deleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := getIdFromRequest(r)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	if err := h.service.Delete(id); err != nil {
		http.Error(w, "failed to delete user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
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
