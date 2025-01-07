package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ksusonic/kanban/internal/models"
	"github.com/ksusonic/kanban/internal/server/api"
)

const validationErrorMessage = "telegram query data invalid"

func (ctrl *Controller) TelegramCallback(c *gin.Context) {
	ctx := c.Request.Context()

	var callbackData models.TelegramCallback
	if err := c.ShouldBindQuery(&callbackData); err != nil {
		ctrl.log.WarnContext(ctx, "binding telegram callback", err)
		c.JSON(http.StatusBadRequest, api.ErrorResponse{
			Error: api.ErrorResponseValidationError(errors.New(validationErrorMessage)),
		})
		return
	}

	if !ctrl.authModule.ValidateTelegramCallbackData(c.Request.URL.Query()) {
		ctrl.log.WarnContext(ctx, validationErrorMessage, "query", c.Request.URL.RawQuery)
		c.JSON(http.StatusBadRequest, api.ErrorResponse{
			Error: api.ErrorResponseValidationError(errors.New(validationErrorMessage)),
		})
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

func (ctrl *Controller) handleTGUser(ctx context.Context, callbackData *models.TelegramCallback) (int, error) {
	// Check if already registered
	user, err := ctrl.userRepo.GetByTelegramID(ctx, callbackData.ID)
	if err == nil {
		return user.ID, nil
	}

	// Register or error
	if errors.Is(err, models.ErrNotFound) {
		ctrl.log.InfoContext(ctx, "new user", "username", callbackData.Username)

		var userID int
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
