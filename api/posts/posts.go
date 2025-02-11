package posts

import (
	"errors"
	"github.com/go-chi/chi/v5"
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

const label = "post"

type Handler struct {
	logger   *slog.Logger
	query    *database.Queries
	validate *validator.Validate
}

type postRequest struct {
	Title   string `json:"title,required"`
	Content string `json:"content,required"`
}

func RegisterRoutes(r chi.Router, handler *Handler) {
	r.Route("/post", func(r chi.Router) {
		r.Get("/{id}", handler.GetPost)
		r.With(authmiddleware.JWTAuth).Post("/", handler.CreatePost)
		r.With(authmiddleware.JWTAuth).Delete("/{id}", handler.DeletePost)
		r.With(authmiddleware.JWTAuth).Put("/{id}", handler.UpdatePost)
	})

	r.Get("/posts", handler.GetPosts)
}

func NewPostsHandler(log *slog.Logger, db *database.Queries) *Handler {
	return &Handler{
		logger:   log,
		query:    db,
		validate: validator.New(),
	}
}

func (h *Handler) GetPost(w http.ResponseWriter, r *http.Request) {

	const op = "posts.GetPost"

	idAlias := chi.URLParam(r, "id")

	id, err := uuid.Parse(idAlias)

	if err != nil {
		h.logger.Warn(op, "failed to parse id as a valid uuid", sl.Err(err), slog.String("id", idAlias))
		json.WriteJSON(w, http.StatusBadRequest, response.BadRequest("passed invalid id"))
		return
	}

	post, err := h.query.GetPost(r.Context(), id)

	if err != nil {
		errD := sqlhelpers.GetDBError(err, label)
		h.logger.Warn(op, "failed to get post", sl.Err(err))
		json.WriteJSON(w, errD.StatusCode, errD)
		return
	}

	json.WriteJSON(w, http.StatusOK, response.OkWData(post))
}

func (h *Handler) GetPosts(w http.ResponseWriter, r *http.Request) {
	const op = "posts.GetPosts"

	userId := uuid.NullUUID{Valid: false}

	if possibleId, _, err := authmiddleware.Identify(r.Context(), w, h.logger, op); err == nil {
		userId = uuid.NullUUID{UUID: possibleId, Valid: true}
	}

	posts, err := h.query.GetPostsWithLikes(r.Context(), userId.UUID)

	if err != nil {
		errD := sqlhelpers.GetDBError(err, label)
		h.logger.Warn(op, "failed to get post", sl.Err(err))
		json.WriteJSON(w, errD.StatusCode, errD)
		return
	}

	if len(posts) == 0 {
		json.WriteJSON(w, http.StatusOK, response.OkWData([]string{}))
		return
	}

	json.WriteJSON(w, http.StatusOK, response.OkWData(posts))
}

func (h *Handler) DeletePost(w http.ResponseWriter, r *http.Request) {
	const op = "posts.DeletePost"

	idAlias := chi.URLParam(r, "id")

	postId, err := uuid.Parse(idAlias)

	if err != nil {
		h.logger.Warn(op, "failed to parse id as a valid uuid", sl.Err(err), slog.String("id", idAlias))
		json.WriteJSON(w, http.StatusBadRequest, response.BadRequest("passed invalid id"))
		return
	}

	authorId, errD, err := authmiddleware.Identify(r.Context(), w, h.logger, op)

	if err != nil {
		json.WriteJSON(w, errD.StatusCode, errD)
		return
	}

	post, err := h.query.GetPost(r.Context(), postId)

	if err != nil {
		h.logger.Warn("Failed to get post", sl.Err(err))
		errD := sqlhelpers.GetDBError(err, label)
		json.WriteJSON(w, errD.StatusCode, errD)
		return
	}

	if post.AuthorID != authorId {
		json.WriteJSON(w, http.StatusForbidden, response.Forbidden("You cannot delete this post"))
		return
	}

	err = h.query.DeletePost(r.Context(), post.ID)

	if err != nil {
		h.logger.Warn("Failed to delete post", sl.Err(err))
		errD := sqlhelpers.GetDBError(err, label)
		json.WriteJSON(w, errD.StatusCode, errD)
		return
	}

	json.WriteJSON(w, http.StatusOK, response.OkWMsg("Post deleted"))
}

func (h *Handler) CreatePost(w http.ResponseWriter, r *http.Request) {

	const op = "posts.CreatePost"

	var req postRequest

	authorId, errD, err := authmiddleware.Identify(r.Context(), w, h.logger, op)

	if err != nil {
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

	post, err := h.query.CreatePost(r.Context(), database.CreatePostParams{
		ID:        uuid.New(),
		AuthorID:  authorId,
		Title:     req.Title,
		Content:   req.Content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	if err != nil {
		h.logger.Warn("Failed to create post", sl.Err(err))
		errD := sqlhelpers.GetDBError(err, label)
		json.WriteJSON(w, errD.StatusCode, errD)
		return
	}

	json.WriteJSON(w, http.StatusCreated, response.OkWDataAMsg(post, "Post created successfully"))
}

func (h *Handler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	const op = "posts.CreatePost"

	idAlias := chi.URLParam(r, "id")

	postId, err := uuid.Parse(idAlias)

	if err != nil {
		h.logger.Warn(op, "failed to parse id as a valid uuid", sl.Err(err), slog.String("id", idAlias))
		json.WriteJSON(w, http.StatusBadRequest, response.BadRequest("passed invalid id"))
		return
	}

	var req postRequest

	authorId, errD, err := authmiddleware.Identify(r.Context(), w, h.logger, op)

	if err != nil {
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

	post, err := h.query.GetPost(r.Context(), postId)

	if err != nil {
		h.logger.Warn("Failed to get post", sl.Err(err))
		errD := sqlhelpers.GetDBError(err, label)
		json.WriteJSON(w, errD.StatusCode, errD)
		return
	}

	if post.AuthorID != authorId {
		json.WriteJSON(w, http.StatusForbidden, response.Forbidden("You cannot update this post"))
		return
	}

	updatedP, err := h.query.UpdatePost(r.Context(), database.UpdatePostParams{
		ID:        uuid.New(),
		Title:     req.Title,
		Content:   req.Content,
		UpdatedAt: time.Now(),
	})

	if err != nil {
		h.logger.Warn("Failed to update post", sl.Err(err))
		errD := sqlhelpers.GetDBError(err, label)
		json.WriteJSON(w, http.StatusInternalServerError, errD)
		return
	}

	json.WriteJSON(w, http.StatusOK, response.OkWDataAMsg(updatedP, "Post updated successfully"))
}
