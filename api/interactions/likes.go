package interactions

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"log/slog"
	authmiddleware "money-manager/api/middlewares/auth"
	"money-manager/internal/database"
	"money-manager/internal/lib/http/json"
	"money-manager/internal/lib/http/response"
	"money-manager/internal/lib/logger/sl"
	"money-manager/internal/lib/sql/sqlhelpers"
	"net/http"
	"time"
)

const likeLabel = "like"

type likePostRequest struct {
	EntityID   string `json:"entity_id" validate:"required,uuid"`
	EntityType string `json:"entity_type" validate:"required,oneof=post comment"`
}

func (h *Handler) LikeEntity(w http.ResponseWriter, r *http.Request) {
	const op = "interactions.LikeEntity"

	var req likePostRequest

	if details, err := json.DecodeJSONBody(w, r, &req); err != nil {
		h.logger.Warn("Invalid JSON body", slog.String("op", op), sl.Err(err))
		json.WriteJSON(w, details.StatusCode, details)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			h.logger.Warn("Validation failed", slog.String("op", op), sl.Err(err))
			json.WriteJSON(w, http.StatusBadRequest, response.InvalidInput(validationErrors))
			return
		}

		h.logger.Warn("Invalid input data", slog.String("op", op), sl.Err(err))
		response.BadRequest("invalid input data")
		return
	}

	entityID, err := uuid.Parse(req.EntityID)

	if err != nil {
		h.logger.Warn("Invalid entity id", slog.String("op", op), sl.Err(err))
		errD := response.BadRequest("invalid entity id")
		json.WriteJSON(w, errD.StatusCode, errD)
		return
	}

	currentUserId, errD, err := authmiddleware.Identify(r, w, h.logger, op)

	if err != nil {
		json.WriteJSON(w, errD.StatusCode, errD)
		return
	}

	likedRows, err := h.query.LikeEntity(r.Context(), database.LikeEntityParams{
		ID:         uuid.New(),
		UserID:     currentUserId,
		EntityID:   entityID,
		EntityType: req.EntityType,
		CreatedAt:  time.Now(),
	})

	if err != nil {
		h.logger.Warn("Like failed", slog.String("op", op), sl.Err(err))
		errD = sqlhelpers.GetDBError(err, likeLabel)
		json.WriteJSON(w, errD.StatusCode, errD)
		return
	}

	if likedRows == 0 {
		h.logger.Warn("Attempt to like non-existent post/comment", slog.String("op", op))
		errD = response.NotFound(fmt.Sprintf("%s does not exist", req.EntityType))
		json.WriteJSON(w, errD.StatusCode, errD)
		return
	}

	json.WriteJSON(w, http.StatusOK, response.OkWMsg("Successfully liked"))
}

func (h *Handler) UnlikeEntity(w http.ResponseWriter, r *http.Request) {
	const op = "interactions.UnlikePost"

	var req likePostRequest

	if details, err := json.DecodeJSONBody(w, r, &req); err != nil {
		h.logger.Warn("Invalid JSON body", slog.String("op", op), sl.Err(err))
		json.WriteJSON(w, details.StatusCode, details)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			h.logger.Warn("Validation failed", slog.String("op", op), sl.Err(err))
			json.WriteJSON(w, http.StatusBadRequest, response.InvalidInput(validationErrors))
			return
		}

		h.logger.Warn("Invalid input data", slog.String("op", op), sl.Err(err))
		response.BadRequest("invalid input data")
		return
	}

	entityID, err := uuid.Parse(req.EntityID)

	if err != nil {
		h.logger.Warn("Invalid entity id", slog.String("op", op), sl.Err(err))
		errD := response.BadRequest("invalid entity id")
		json.WriteJSON(w, errD.StatusCode, errD)
		return
	}

	currentUserId, errD, err := authmiddleware.Identify(r, w, h.logger, op)

	if err != nil {
		json.WriteJSON(w, errD.StatusCode, errD)
		return
	}

	unlikedRows, err := h.query.UnlikeEntity(r.Context(), database.UnlikeEntityParams{
		UserID:     currentUserId,
		EntityID:   entityID,
		EntityType: req.EntityType,
	})

	if err != nil {
		h.logger.Warn("Unlike failed", slog.String("op", op), sl.Err(err))
		errD = sqlhelpers.GetDBError(err, likeLabel)
		json.WriteJSON(w, errD.StatusCode, errD)
		return
	}

	if unlikedRows == 0 {
		errD = response.NotFound(fmt.Sprintf("%s like does not exist", req.EntityType))
		json.WriteJSON(w, errD.StatusCode, errD)
		return
	}

	json.WriteJSON(w, http.StatusOK, response.OkWMsg("Successfully unliked"))
}
