package board

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/ksusonic/kanban/internal/server/api"
	"github.com/ksusonic/kanban/internal/server/requestctx"
)

func (ctrl *Controller) AvailableBoards(c *gin.Context) {
	boards, err := ctrl.feature.AvailableBoards(c.Request.Context(), requestctx.MustGetUser(c).UserID)
	if err != nil {
		return
	}

	response := api.AvailableBoardsResponse{
		Boards: make([]api.BoardShortInfo, 0, len(boards)),
	}
	for i := range boards {
		response.Boards = append(response.Boards, api.BoardShortInfo{
			ID:        boards[i].ID,
			Slug:      boards[i].Slug,
			Name:      boards[i].Name,
			CreatedAt: boards[i].CreatedAt.Format(time.DateOnly),
		})
	}

	c.JSON(http.StatusOK, response)
}
