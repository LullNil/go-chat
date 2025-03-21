package auth

import (
	"encoding/json"
	"fmt"
	"net/http"

	ssogrpc "go-chat/internal/clients/sso/grpc"
)

// enableCors sets Access-Control-Allow-Origin header for CORS.
func enableCors(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, appID")
}

type Handler struct {
	SSOClient *ssogrpc.Client
}

func NewAuthHandler(ssoClient *ssogrpc.Client) *Handler {
	return &Handler{SSOClient: ssoClient}
}

// RegisterHandler handles requests to register a user.
func (h *Handler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(w)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	userID, err := h.SSOClient.Register(r.Context(), req.Email, req.Password)
	if err != nil {
		fmt.Printf("gRPC register error: %v\n", err)

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"userId": userID,
	})
}

// LoginHandler handles requests to login a user.
func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(w)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		AppID    int32  `json:"app_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	token, err := h.SSOClient.Login(r.Context(), req.Email, req.Password,
		req.AppID)
	if err != nil {
		fmt.Printf("gRPC login error: %v\n", err)

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"token": token,
	})
}
