package rest

import (
	"encoding/json"
	"fmt"
	"github.com/dankru/Commissions_simple/internal/domain"
	"github.com/gorilla/mux"
	"io"
	"net/http"
)

func (h *Handler) initUserRoutes(router *mux.Router) {
	users := router.PathPrefix("/users").Subrouter()
	{
		users.Use(h.authMiddleware)
		users.HandleFunc("", h.getUsers).Methods(http.MethodGet)
		users.HandleFunc("/{id:[0-9]+}", h.getUserById).Methods(http.MethodGet)
		users.HandleFunc("/{id:[0-9]+}", h.replaceUser).Methods(http.MethodPut)
		users.HandleFunc("/{id:[0-9]+}", h.updateUser).Methods(http.MethodPatch)
		users.HandleFunc("/{id:[0-9]+}", h.deleteUser).Methods(http.MethodDelete)
	}
}

func (h *Handler) getUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.userService.GetAll()
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

	user, err := h.userService.GetById(id)
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

	if err := h.userService.Replace(id, user); err != nil {
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

	if err := h.userService.Update(id, userInp); err != nil {
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

	if err := h.userService.Delete(id); err != nil {
		http.Error(w, "failed to delete user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
