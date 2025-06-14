package user

import (
	"encoding/json"
	"net/http"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		// не переданы поля согласно структуре
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	_, err := h.svc.Register(creds.Username, creds.Password)
	if err != nil {
		// ошибка регистрации
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("registered"))
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		// не переданы поля согласно структуре
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	accessToken, refreshToken, err := h.svc.Login(creds.Username, creds.Password)
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
