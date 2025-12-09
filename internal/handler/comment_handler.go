package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"khalif-comment/internal/domain"
	"khalif-comment/pkg/utils"

)

type CommentHandler struct {
	useCase domain.CommentUseCase
}

func NewCommentHandler(u domain.CommentUseCase) *CommentHandler {
	return &CommentHandler{useCase: u}
}

// --- DTOs ---

type CreateCommentRequest struct {
	StoryID string `json:"story_id" binding:"required"`
	Content string `json:"content" binding:"required"`
}

type UpdateCommentRequest struct {
	Content string `json:"content" binding:"required"`
}

// --- HANDLERS ---

// CreateComment godoc
// @Summary      Post a comment
// @Description  Create a new comment for a story
// @Tags         comments
// @Accept       json
// @Produce      json
// @Param        request body CreateCommentRequest true "Comment Data"
// @Success      201  {object}  domain.Comment
// @Failure      400  {object}  utils.APIResponse
// @Failure      500  {object}  utils.APIResponse
// @Router       /comments [post]
// @Security     BearerAuth
func (h *CommentHandler) Create(c *gin.Context) {
	var req CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	userID := c.GetString("user_id") // Ambil dari Middleware Auth
	
	res, err := h.useCase.Create(c.Request.Context(), req.StoryID, userID, req.Content)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusCreated, res)
}

// GetCommentsByStory godoc
// @Summary      Get comments by story
// @Description  Retrieve all comments for a specific story
// @Tags         comments
// @Produce      json
// @Param        story_id query     string  true  "Story UUID"
// @Success      200  {array}   domain.Comment
// @Failure      400  {object}  utils.APIResponse
// @Router       /comments [get]
func (h *CommentHandler) GetByStory(c *gin.Context) {
	storyID := c.Query("story_id")
	if storyID == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "story_id query parameter is required")
		return
	}

	res, err := h.useCase.GetByStoryUUID(c.Request.Context(), storyID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, res)
}

// UpdateComment godoc
// @Summary      Update a comment
// @Description  Update comment content (Owner only)
// @Tags         comments
// @Accept       json
// @Produce      json
// @Param        id       path      int                   true  "Comment ID"
// @Param        request  body      UpdateCommentRequest  true  "Update Data"
// @Success      200  {object}  domain.Comment
// @Failure      400  {object}  utils.APIResponse
// @Failure      403  {object}  utils.APIResponse "Unauthorized Action"
// @Failure      404  {object}  utils.APIResponse
// @Failure      500  {object}  utils.APIResponse
// @Router       /comments/{id} [put]
// @Security     BearerAuth
func (h *CommentHandler) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid comment id")
		return
	}

	var req UpdateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	userID := c.GetString("user_id")

	res, err := h.useCase.Update(c.Request.Context(), uint(id), userID, req.Content)
	if err != nil {
		if errors.Is(err, domain.ErrCommentNotFound) {
			utils.ErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}
		if errors.Is(err, domain.ErrUnauthorizedAction) {
			utils.ErrorResponse(c, http.StatusForbidden, err.Error()) // 403 Forbidden
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, res)
}

// DeleteComment godoc
// @Summary      Delete a comment
// @Description  Delete a comment by ID (Owner only)
// @Tags         comments
// @Produce      json
// @Param        id   path      int     true  "Comment ID"
// @Success      200  {object}  utils.APIResponse
// @Failure      403  {object}  utils.APIResponse "Unauthorized Action"
// @Failure      404  {object}  utils.APIResponse
// @Failure      500  {object}  utils.APIResponse
// @Router       /comments/{id} [delete]
// @Security     BearerAuth
func (h *CommentHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid comment id")
		return
	}

	userID := c.GetString("user_id")

	if err := h.useCase.Delete(c.Request.Context(), uint(id), userID); err != nil {
		if errors.Is(err, domain.ErrCommentNotFound) {
			utils.ErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}
		if errors.Is(err, domain.ErrUnauthorizedAction) {
			utils.ErrorResponse(c, http.StatusForbidden, err.Error()) // 403 Forbidden
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SuccessMessage(c, http.StatusOK, "comment deleted")
}