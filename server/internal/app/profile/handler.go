package profile

import (
	"encoding/json"
	"errors"
	"go-chat/internal/lib/logger/sl"
	"go-chat/internal/lib/response"
	"log/slog"
	"net/http"

	mw "go-chat/internal/app/middleware"
	userservice "go-chat/internal/service/user"
)

// Handler processes requests related to the user profile
type Handler struct {
	userService userservice.UserService
	log         *slog.Logger
}

// NewHandler creates a new instance of Profile Handler
func NewHandler(userService userservice.UserService, log *slog.Logger) *Handler {
	return &Handler{
		userService: userService,
		log:         log,
	}
}

// HandleGetProfile handles GET requests to get the current user's profile
func (h *Handler) GetProfile(w http.ResponseWriter, r *http.Request) {
	const op = "handler.profile.GetProfile"
	log := h.log.With(slog.String("op", op))

	// Get userID from context
	userID, ok := mw.GetUserIDFromContext(r.Context())
	if !ok {
		log.Error("userID not found in context")
		response.RespondError(w, http.StatusInternalServerError,
			"Internal server error")
		return
	}

	log.Info("getting full profile", slog.Int64("user_id", userID))

	// Get user profile from repo
	fullProfile, err := h.userService.GetFullUserProfile(r.Context(), userID)
	if err != nil {
		if errors.Is(err, userservice.ErrUserNotFound) || errors.Is(err,
			userservice.ErrProfileNotFound) {
			log.Warn("full profile not found", sl.Err(err),
				slog.Int64("user_id", userID))
			response.RespondError(w, http.StatusNotFound, "Profile data not found")
			return
		}
		log.Error("failed to get full profile from service", sl.Err(err),
			slog.Int64("user_id", userID))
		response.RespondError(w, http.StatusInternalServerError,
			"Failed to get profile")
		return
	}

	log.Info("full profile retrieved successfully", slog.Int64("user_id", userID))
	response.RespondJSON(w, http.StatusOK, fullProfile)
}

// UpdateProfileRequest defines the body of the request to update the profile
type UpdateProfileRequest struct {
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	AvatarURL *string `json:"avatar_url"`
}

func (h *Handler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	const op = "handler.profile.UpdateProfile"
	log := h.log.With(slog.String("op", op))

	// Get userID from context
	userID, ok := mw.GetUserIDFromContext(r.Context())
	if !ok {
		log.Error("userID not found in context")
		response.RespondError(w, http.StatusInternalServerError,
			"Internal server error")
		return
	}

	log.Info("updating profile", slog.Int64("user_id", userID))

	var req UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error("failed to decode request body", sl.Err(err))
		response.RespondError(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	updateDTO := &userservice.UpdateProfileRequestDTO{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		AvatarURL: req.AvatarURL,
	}

	updatedProfile, err := h.userService.UpdateUserProfile(r.Context(),
		userID, updateDTO)
	if err != nil {
		log.Error("failed to update profile via service", sl.Err(err),
			slog.Int64("user_id", userID))
		response.RespondError(w, http.StatusInternalServerError,
			"Failed to update profile")
		return
	}

	log.Info("profile updated successfully", slog.Int64("user_id", userID))

	if updatedProfile != nil {
		response.RespondJSON(w, http.StatusOK, updatedProfile)
	} else {
		response.RespondJSON(w, http.StatusOK, map[string]string{"status": "updated"})
	}
}

// TODO: Add a handler for GetPublicProfileByUsername that will
// extract {username} from the URL (using chi.URLParam) and call
// h.userService.GetPublicProfileByUsername
