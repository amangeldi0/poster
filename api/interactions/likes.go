package interactions

import (
	"errors"
	"github.com/go-chi/chi/v5"
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

func (h *Handler) isValidUUIDParam(r *http.Request) (uuid.UUID, error) {
	paramID := chi.URLParam(r, "id")
	invalidID := errors.New("invalid id")

	if paramID == "" {
		return uuid.Nil, invalidID
	}

	id, err := uuid.Parse(paramID)

	if err != nil {
		return uuid.Nil, invalidID
	}

	return id, nil
}

func (h *Handler) LikeComment(w http.ResponseWriter, r *http.Request) {
	const op = "interactions.LikeComment"

	commentID, err := h.isValidUUIDParam(r)

	if err != nil {
		h.logger.Warn("invalid comment id", slog.String("op", op), sl.Err(err))
		errD := response.BadRequest(err.Error())
		json.WriteJSON(w, errD.StatusCode, errD)
		return
	}

	currentUserId, errD, err := authmiddleware.Identify(r, w, h.logger, op)

	if err != nil {
		json.WriteJSON(w, errD.StatusCode, errD)
		return
	}

	if ok := h.isCommentExist(r.Context(), commentID); !ok {
		h.logger.Warn("attempt to like non-existent comment", slog.String("op", op))
		errD = response.NotFound(errCommentNotFound.Error())
		json.WriteJSON(w, errD.StatusCode, errD)
		return
	}

	likedRows, err := h.query.LikeComment(r.Context(), database.LikeCommentParams{
		ID:        uuid.New(),
		UserID:    currentUserId,
		CommentID: commentID,
		CreatedAt: time.Now(),
	})

	if err != nil {
		h.logger.Warn("like failed", slog.String("op", op), sl.Err(err))
		errD = sqlhelpers.GetDBError(err, commentLabel)
		json.WriteJSON(w, errD.StatusCode, errD)
		return
	}

	if likedRows == 0 {
		h.logger.Warn("attempt to like non-existent comment", slog.String("op", op))
		errD = response.NotFound(errCommentAttachmentNotFound.Error())
		json.WriteJSON(w, errD.StatusCode, errD)
		return
	}

	json.WriteJSON(w, http.StatusOK, response.OkWMsg("comment successfully liked"))
}

func (h *Handler) UnlikeComment(w http.ResponseWriter, r *http.Request) {
	const op = "interactions.UnlikeComment"

	commentID, err := h.isValidUUIDParam(r)

	if err != nil {
		h.logger.Warn("invalid comment id", slog.String("op", op), sl.Err(err))
		errD := response.BadRequest(err.Error())
		json.WriteJSON(w, errD.StatusCode, errD)
		return
	}

	currentUserId, errD, err := authmiddleware.Identify(r, w, h.logger, op)

	if err != nil {
		json.WriteJSON(w, errD.StatusCode, errD)
		return
	}

	if ok := h.isCommentExist(r.Context(), commentID); !ok {
		h.logger.Warn("attempt to like non-existent comment", slog.String("op", op))
		errD = response.NotFound(errCommentNotFound.Error())
		json.WriteJSON(w, errD.StatusCode, errD)
		return
	}

	likedRows, err := h.query.UnlikeComment(r.Context(), database.UnlikeCommentParams{
		UserID:    currentUserId,
		CommentID: commentID,
	})

	if err != nil {
		h.logger.Warn("unlike failed", slog.String("op", op), sl.Err(err))
		errD = sqlhelpers.GetDBError(err, commentLabel)
		json.WriteJSON(w, errD.StatusCode, errD)
		return
	}

	if likedRows == 0 {
		h.logger.Warn("attempt to like non-existent comment", slog.String("op", op))
		errD = response.NotFound(errCommentAttachmentNotFound.Error())
		json.WriteJSON(w, errD.StatusCode, errD)
		return
	}

	json.WriteJSON(w, http.StatusOK, response.OkWMsg("comment successfully unliked"))
}

func (h *Handler) LikePost(w http.ResponseWriter, r *http.Request) {
	const op = "interactions.LikePost"

	postID, err := h.isValidUUIDParam(r)

	if err != nil {
		h.logger.Warn("invalid comment id", slog.String("op", op), sl.Err(err))
		errD := response.BadRequest(err.Error())
		json.WriteJSON(w, errD.StatusCode, errD)
		return
	}

	currentUserId, errD, err := authmiddleware.Identify(r, w, h.logger, op)

	if err != nil {
		json.WriteJSON(w, errD.StatusCode, errD)
		return
	}

	if ok := h.isPostExist(r.Context(), postID); !ok {
		h.logger.Warn("attempt to like non-existent post", slog.String("op", op))
		errD = response.NotFound(errPostNotFound.Error())
		json.WriteJSON(w, errD.StatusCode, errD)
		return
	}

	likedRows, err := h.query.LikePost(r.Context(), database.LikePostParams{
		ID:        uuid.New(),
		UserID:    currentUserId,
		PostID:    postID,
		CreatedAt: time.Now(),
	})

	if err != nil {
		h.logger.Warn("like failed", slog.String("op", op), sl.Err(err))
		errD = sqlhelpers.GetDBError(err, postLabel)
		json.WriteJSON(w, errD.StatusCode, errD)
		return
	}

	if likedRows == 0 {
		h.logger.Warn("attempt to like non-existent post", slog.String("op", op))
		errD = response.NotFound(errPostAttachmentNotFound.Error())
		json.WriteJSON(w, errD.StatusCode, errD)
		return
	}

	json.WriteJSON(w, http.StatusOK, response.OkWMsg("post successfully liked"))
}

func (h *Handler) UnlikePost(w http.ResponseWriter, r *http.Request) {
	const op = "interactions.UnlikePost"

	postID, err := h.isValidUUIDParam(r)

	if err != nil {
		h.logger.Warn("invalid comment id", slog.String("op", op), sl.Err(err))
		errD := response.BadRequest(err.Error())
		json.WriteJSON(w, errD.StatusCode, errD)
		return
	}

	currentUserId, errD, err := authmiddleware.Identify(r, w, h.logger, op)

	if err != nil {
		json.WriteJSON(w, errD.StatusCode, errD)
		return
	}

	if ok := h.isPostExist(r.Context(), postID); !ok {
		h.logger.Warn("attempt to like non-existent post", slog.String("op", op))
		errD = response.NotFound(errPostNotFound.Error())
		json.WriteJSON(w, errD.StatusCode, errD)
		return
	}

	likedRows, err := h.query.UnlikePost(r.Context(), database.UnlikePostParams{
		UserID: currentUserId,
		PostID: postID,
	})

	if err != nil {
		h.logger.Warn("unlike failed", slog.String("op", op), sl.Err(err))
		errD = sqlhelpers.GetDBError(err, postLabel)
		json.WriteJSON(w, errD.StatusCode, errD)
		return
	}

	if likedRows == 0 {
		h.logger.Warn("attempt to like non-existent post", slog.String("op", op))
		errD = response.NotFound(errPostAttachmentNotFound.Error())
		json.WriteJSON(w, errD.StatusCode, errD)
		return
	}

	json.WriteJSON(w, http.StatusOK, response.OkWMsg("post successfully unliked"))
}
