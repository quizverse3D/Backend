package authgateway

import (
	"encoding/json"
	"net/http"
	"strings"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	UserName string `json:"username"` // только для передачи по шине Users
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var creds Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		// не переданы поля согласно структуре
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	if creds.Email == "" || creds.Password == "" || creds.UserName == "" {
		http.Error(w, "email, password, username must be provided", http.StatusBadRequest)
		return
	}

	_, err := h.svc.Register(creds.Email, creds.Password, creds.UserName)
	if err != nil {
		// ошибка регистрации
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("registered"))
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var creds Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		// не переданы поля согласно структуре
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	if creds.Email == "" || creds.Password == "" {
		http.Error(w, "email and password must be provided", http.StatusBadRequest)
		return
	}

	accessToken, refreshToken, err := h.svc.Login(creds.Email, creds.Password)
	if err != nil {
		// ошибка авторизации
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HttpOnly: true,
		Path:     "/",
		MaxAge:   7 * 24 * 60 * 60, // 7 дней
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(accessToken))
}

func (h *Handler) ValidateToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var payload struct {
		AccessToken string `json:"accessToken"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	userID, err := h.svc.ValidateAccessToken(payload.AccessToken)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	response := map[string]string{"userId": userID}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) RefreshAccessToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		http.Error(w, "refresh token not provided", http.StatusUnauthorized)
		return
	}

	newAccessToken, err := h.svc.RefreshAccessToken(cookie.Value)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(newAccessToken))
}

func (h *Handler) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	token := r.Header.Get("Authorization")
	if token == "" {
		http.Error(w, "missing token", http.StatusUnauthorized)
		return
	}

	userUuid, err := ValidateAccessToken(strings.TrimPrefix(token, "Bearer "))
	if err != nil {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}

	var payload struct {
		NewPassword string `json:"new_password"`
		OldPassword string `json:"old_password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	if payload.OldPassword == "" || payload.NewPassword == "" {
		http.Error(w, "old_password and new_password must be provided", http.StatusBadRequest)
		return
	}

	err = h.svc.UpdatePassword(userUuid, payload.NewPassword, payload.OldPassword)
	if err == ErrInvalidPassword {
		http.Error(w, "old_password is invalid", http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
