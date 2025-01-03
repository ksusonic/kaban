package auth

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ksusonic/kanban/internal/models"
	"github.com/ksusonic/kanban/internal/telegram"
)

type TelegramCallbackData struct {
	ID        int64   `form:"id" binding:"required"`
	FirstName string  `form:"first_name" binding:"required"`
	Username  string  `form:"username" binding:"required"`
	PhotoURL  *string `form:"photo_url"`
	AuthDate  int64   `form:"auth_date"`
	Hash      string  `form:"hash" binding:"required"`
	Next      string  `form:"next"`
}

func (ctrl *Controller) TelegramCallback(c *gin.Context) {
	ctx := c.Request.Context()

	var callbackData TelegramCallbackData
	if err := c.ShouldBindQuery(&callbackData); err != nil {
		ctrl.log.WarnContext(ctx, "binding telegram callback", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "validation"})
		return
	}

	if !telegram.ValidateTelegramCallbackData(c.Request.URL.Query(), ctrl.botCfg.Token) {
		ctrl.log.WarnContext(ctx, "telegram query data invalid", "query", c.Request.URL.RawQuery)
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad query"})
		return
	}

	ctrl.log.DebugContext(ctx, "telegram callback", "callback", callbackData)

	userID, err := ctrl.userRepo.GetUserIDByTelegramID(ctx, callbackData.ID)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrNotFound):
			ctrl.log.InfoContext(ctx, "new user", "username", callbackData.Username)

			_, err = ctrl.userRepo.AddTelegramUser(
				ctx,
				callbackData.Username,
				callbackData.ID,
				callbackData.FirstName,
				callbackData.PhotoURL,
			)
			if err != nil {
				ctrl.log.ErrorContext(ctx, "add new tg user", err, "user", callbackData)
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
		default:
			ctrl.log.ErrorContext(ctx, "get user by telegram", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}

	ctrl.tgAuth(c, userID)

	c.JSON(http.StatusOK, gin.H{"user_id": userID})
}

func (ctrl *Controller) tgAuth(_ *gin.Context, _ int) {
	panic("not implemented")
}
