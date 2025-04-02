package auth

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	ssogrpc "go-chat/internal/clients/sso/grpc"
	"go-chat/internal/lib/jwtutil"
	"go-chat/internal/lib/logger"
	"go-chat/internal/lib/logger/sl"
	"go-chat/internal/storage/redis"
)

type Handler struct {
	SSOClient   *ssogrpc.Client
	RedisClient *redis.Client
}

func NewAuthHandler(ssoClient *ssogrpc.Client, redisClient *redis.Client) *Handler {
	return &Handler{SSOClient: ssoClient, RedisClient: redisClient}
}

// RegisterHandler handles requests to register a user.
func (h *Handler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Username string `json:"username"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	_, err := h.SSOClient.Register(r.Context(), req.Email, req.Password,
		req.Username)
	if err != nil {
		logger.GetLogger().Error("gRPC register error", sl.Err(err))
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"status": "register_success",
	})
}

// LoginHandler handles requests to login a user.
func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Username string `json:"username"`
		AppID    int32  `json:"app_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	token, err := h.SSOClient.Login(r.Context(), req.Email, req.Password,
		req.Username, req.AppID)
	if err != nil {
		logger.GetLogger().Error("gRPC login error", sl.Err(err))
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	cookie := http.Cookie{
		Name:     "token",
		Value:    token,
		HttpOnly: true,
		Secure:   false, // TODO: Change to true in production
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
	}
	http.SetCookie(w, &cookie)

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"status": "login_success",
	})
}

// GetUserHandler handles requests to get user data.
func (h *Handler) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	// Get token from cookie.
	cookie, err := r.Cookie("token")
	if err != nil {
		respondError(w, http.StatusUnauthorized, err.Error())
		return
	}

	// Verify token and get user id.
	userID, err := jwtutil.VerifyToken(r.Context(), cookie.Value)
	if err != nil {
		code := http.StatusInternalServerError
		if err.Error() == "token expired" {
			code = http.StatusUnauthorized
		}
		respondError(w, code, err.Error())
		return
	}

	// Try to get user data from Redis.
	userData, err := h.RedisClient.GetCache(r.Context(), userID)
	switch {
	case err != nil:
		logger.GetLogger().Error("failed to get user data from redis", sl.Err(err))
	case userData == nil:
		logger.GetLogger().Info("cache miss", slog.Int("userID", int(userID)))
	default:
		logger.GetLogger().Info("cache hit", slog.Int("userID", int(userID)))
		respondJSON(w, http.StatusOK, userData)
		return
	}

	// If there is no data in Redis or an error occurred while retrieving,
	// turn to gRPC.
	userData, err = h.SSOClient.GetUser(r.Context(), userID)
	if err != nil {
		logger.GetLogger().Error("gRPC get user error", sl.Err(err))
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Cache the received data in Redis with a TTL of 5 minutes.
	if err := h.RedisClient.SetCache(r.Context(), userID,
		userData, 5*time.Minute); err != nil {
		logger.GetLogger().Error("failed to cache user data", sl.Err(err))
	} else {
		logger.GetLogger().Info("user data cached successfully", slog.Int(
			"userID", int(userID)))
	}

	respondJSON(w, http.StatusOK, userData)
}

// LogoutHandler handles requests to log out a user by expiring the
// authentication cookie.
func (h *Handler) LogoutHandler(w http.ResponseWriter, r *http.Request) {

	expiredCookie := http.Cookie{
		Name:     "token",
		Value:    "",
		HttpOnly: true,
		Secure:   false, // TODO: Change to true in production
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		Expires:  time.Unix(0, 0),
	}
	http.SetCookie(w, &expiredCookie)
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"status": "logged_out",
	})
}

// respondError sends JSON with an error message.
func respondError(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// respondJSON sends data as JSON.
func respondJSON(w http.ResponseWriter, code int, data interface{}) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}
