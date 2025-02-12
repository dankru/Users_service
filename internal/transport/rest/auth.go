package rest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dankru/Commissions_simple/internal/domain"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
)

func (h *Handler) initAuthRoutes(router *mux.Router) {
	auth := router.PathPrefix("/auth").Subrouter()
	{
		auth.HandleFunc("/sign-up", h.signUp).Methods(http.MethodPost)
		auth.HandleFunc("/sign-in", h.signIn).Methods(http.MethodGet)
		auth.HandleFunc("/refresh", h.refresh).Methods(http.MethodGet)
	}
}

func (h *Handler) signUp(w http.ResponseWriter, r *http.Request) {

	var user domain.UserInput
	user, err := decodeJsonBody[domain.UserInput](r)

	if err = user.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = h.authService.SignUp(user); err != nil {
		http.Error(w, fmt.Sprintf("failed to create user: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) signIn(w http.ResponseWriter, r *http.Request) {
	var signInInput domain.SignInInput

	signInInput, err := decodeJsonBody[domain.SignInInput](r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = signInInput.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	accessToken, refreshToken, err := h.authService.SignIn(signInInput)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := json.Marshal(map[string]string{
		"access_token": accessToken,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Add("Set-Cookie", fmt.Sprintf("refresh-token=%s; HttpOnly", refreshToken))
	w.Header().Add("Content-Type", "application/json")
	w.Write(response)
}

func (h *Handler) refresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh-token")
	if err != nil {
		http.Error(w, "refresh-token cookie not found", http.StatusBadRequest)
		return
	}

	accessToken, refreshToken, err := h.authService.RefreshTokens(cookie.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(map[string]string{
		"token": accessToken,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Set-Cookie", fmt.Sprintf("refresh-token='%s'; HttpOnly", refreshToken))
	w.Header().Add("Content-Type", "application/json")
	w.Write(response)
}

func (h *Handler) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := getTokenFromRequest(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		userId, err := h.authService.ParseToken(r.Context(), token)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), ctxUserId, userId)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func getTokenFromRequest(r *http.Request) (string, error) {
	header := r.Header.Get("Authorization")
	if header == "" {
		return "", errors.New("Empty auth header")
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return "", errors.New("Invalid auth header")
	}

	if len(headerParts[1]) == 0 {
		return "", errors.New("token is empty")
	}

	return headerParts[1], nil
}
