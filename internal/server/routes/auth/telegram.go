package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ksusonic/kanban/internal/auth/telegram"
	"github.com/ksusonic/kanban/internal/models"
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

	userID, err := ctrl.handleTGUser(ctx, &callbackData)
	if err != nil {
		ctrl.log.ErrorContext(ctx, "handle tg user", err, "callback", callbackData)
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	jwt, err := ctrl.authModule.GenerateJWTToken(userID)
	if err != nil {
		ctrl.log.ErrorContext(ctx, "generate token", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":   jwt.Token,
		"expires": jwt.Expires,
	})
}

func (ctrl *Controller) handleTGUser(ctx context.Context, callbackData *TelegramCallbackData) (int, error) {
	// Check if already registered
	userID, err := ctrl.userRepo.GetUserIDByTelegramID(ctx, callbackData.ID)
	if err == nil {
		return userID, nil
	}

	// Register or error
	if errors.Is(err, models.ErrNotFound) {
		ctrl.log.InfoContext(ctx, "new user", "username", callbackData.Username)

		userID, err = ctrl.userRepo.AddTelegramUser(
			ctx,
			callbackData.Username,
			callbackData.ID,
			callbackData.FirstName,
			callbackData.PhotoURL,
		)
		if err != nil {
			return 0, fmt.Errorf("add new: %w", err)
		}

		return userID, nil
	}

	return 0, fmt.Errorf("get userID by telegramID: %w", err)
}
