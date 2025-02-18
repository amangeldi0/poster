package interactions

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
	authmiddleware "poster/api/middlewares/auth"
	"poster/internal/database"
	"poster/internal/lib/http/json"
	"poster/internal/lib/http/response"
	"poster/internal/lib/logger/sl"
	"poster/internal/lib/sql/sqlhelpers"
	"time"
)

type createCommentRequest struct {
	PostId  string `json:"post_id" validate:"required,uuid"`
	Content string `json:"content" validate:"required"`
}

func (h *Handler) Comment(w http.ResponseWriter, r *http.Request) {
	const op = "interactions.comments.Comment"
	var req createCommentRequest

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

	postId, err := uuid.Parse(req.PostId)

	if err != nil {
		h.logger.Warn("Invalid post id", slog.String("op", op), sl.Err(err))
		errD := response.BadRequest("invalid post id")
		json.WriteJSON(w, errD.StatusCode, errD)
		return
	}

	currentUserId, errD, err := authmiddleware.Identify(r, w, h.logger, op)

	if err != nil {
		json.WriteJSON(w, errD.StatusCode, errD)
		return
	}

	comment, err := h.query.CreateComment(r.Context(), database.CreateCommentParams{
		ID:        uuid.New(),
		PostID:    postId,
		UserID:    currentUserId,
		IsEdited:  false,
		Content:   req.Content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	if err != nil {
		h.logger.Warn("Failed to create comment", slog.String("op", op), sl.Err(err))
		errD = sqlhelpers.GetDBError(err, "comment")
		json.WriteJSON(w, errD.StatusCode, errD)
		return
	}

	if comment.ID == uuid.Nil {
		h.logger.Warn("Attempt to comment on non-existent post", slog.String("op", op))
		json.WriteJSON(w, http.StatusNotFound, response.NotFound("Post does not exist"))
		return
	}

	json.WriteJSON(w, http.StatusOK, response.OkWDataAMsg(comment, "comment successfully created"))
}
