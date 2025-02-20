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

type commentRequest struct {
	PostId  string `json:"post_id" validate:"required,uuid"`
	Content string `json:"content" validate:"required"`
}

func (h *Handler) Comment(w http.ResponseWriter, r *http.Request) {
	const op = "interactions.comments.Comment"
	var req commentRequest

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

func (h *Handler) UpdateComment(w http.ResponseWriter, r *http.Request) {
	const op = "interactions.comments.UpdateComment"
	var req commentRequest

	commentID, err := h.isValidUUIDParam(r)

	if err != nil {
		h.logger.Warn("invalid comment id", slog.String("op", op), sl.Err(err))
		errD := response.BadRequest(err.Error())
		json.WriteJSON(w, errD.StatusCode, errD)
		return
	}

	if details, err := json.DecodeJSONBody(w, r, &req); err != nil {
		h.logger.Warn("Invalid JSON body", slog.String("op", op), sl.Err(err))
		json.WriteJSON(w, details.StatusCode, details)
		return
	}

	if err = h.validate.Struct(req); err != nil {
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

	updatedComment, err := h.query.UpdateComment(r.Context(), database.UpdateCommentParams{
		ID:      commentID,
		PostID:  postId,
		UserID:  currentUserId,
		Content: req.Content,
	})

	if err != nil {
		h.logger.Warn("Failed to update comment", slog.String("op", op), sl.Err(err))
		errD = sqlhelpers.GetDBError(err, commentLabel)
		json.WriteJSON(w, errD.StatusCode, errD)
		return
	}

	json.WriteJSON(w, http.StatusOK, response.OkWDataAMsg(updatedComment, "comment successfully updated"))
}

type deleteCommentRequest struct {
	PostId string `json:"post_id" validate:"required,uuid"`
}

func (h *Handler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	const op = "interactions.comments.DeleteComment"
	var req deleteCommentRequest

	commentID, err := h.isValidUUIDParam(r)

	if err != nil {
		h.logger.Warn("invalid comment id", slog.String("op", op), sl.Err(err))
		errD := response.BadRequest(err.Error())
		json.WriteJSON(w, errD.StatusCode, errD)
		return
	}

	if details, err := json.DecodeJSONBody(w, r, &req); err != nil {
		h.logger.Warn("Invalid JSON body", slog.String("op", op), sl.Err(err))
		json.WriteJSON(w, details.StatusCode, details)
		return
	}

	if err = h.validate.Struct(req); err != nil {
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

	deletedRows, err := h.query.DeleteComment(r.Context(), database.DeleteCommentParams{
		ID:     commentID,
		PostID: postId,
		UserID: currentUserId,
	})

	if err != nil {
		h.logger.Warn("Failed to update comment", slog.String("op", op), sl.Err(err))
		errD = sqlhelpers.GetDBError(err, commentLabel)
		json.WriteJSON(w, errD.StatusCode, errD)
		return
	}

	if deletedRows == 0 {
		h.logger.Error("Attempt to delete non-existent comment", slog.String("op", op))
		json.WriteJSON(w, http.StatusNotFound, response.NotFound(errCommentAttachmentNotFound.Error()))
		return
	}

	json.WriteJSON(w, http.StatusOK, response.OkWMsg("comment successfully deleted"))
}
