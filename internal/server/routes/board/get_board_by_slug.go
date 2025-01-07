package board

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/ksusonic/kanban/internal/models"
	"github.com/ksusonic/kanban/internal/server/api"
	"github.com/ksusonic/kanban/internal/server/requestctx"
)

func (ctrl *Controller) GetBoardBySlug(c *gin.Context) {
	ctx := c.Request.Context()

	var request api.GetBoardBySlugRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.JSON(http.StatusBadRequest, api.ErrorResponse{Error: api.ErrorResponseValidationError(err)})
		return
	}

	board, err := ctrl.feature.GetBoardBySlug(ctx, requestctx.MustGetUser(c).UserID, request.Slug)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrNotFound):
			c.JSON(http.StatusNotFound, api.ErrorResponse{Error: api.ErrorResponseBoardNotFound})
		case errors.Is(err, models.ErrNoAccess):
			c.JSON(http.StatusForbidden, api.ErrorResponse{Error: api.ErrorResponseNoAccess})
		default:
			ctrl.log.ErrorContext(ctx, "error getting board by slug", err)
			c.JSON(http.StatusInternalServerError, api.ErrorResponse{Error: api.ErrorResponseInternalServerError})
		}

		return
	}

	c.JSON(http.StatusOK, api.BoardShortInfo{
		ID:        board.ID,
		Slug:      board.Slug,
		Name:      board.Name,
		CreatedAt: board.CreatedAt.Format(time.DateOnly),
	})
}
