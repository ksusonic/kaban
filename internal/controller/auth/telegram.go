package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ksusonic/kanban/internal/telegram"
)

type TelegramCallbackData struct {
	ID        int64  `form:"id"`
	FirstName string `form:"first_name"`
	Username  string `form:"username"`
	PhotoURL  string `form:"photo_url"`
	AuthDate  int64  `form:"auth_date"`
	Hash      string `form:"hash"`
	Next      string `form:"next"`
}

func (ctrl *Controller) TelegramCallback(c *gin.Context) {
	var callbackData TelegramCallbackData
	if err := c.ShouldBindQuery(&callbackData); err != nil {
		ctrl.log.WarnContext(c.Request.Context(), "binding telegram callback", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "validation"})
		return
	}

	if !telegram.ValidateTelegramCallbackData(c.Request.URL.Query(), ctrl.botCfg.Token) {
		ctrl.log.WarnContext(c.Request.Context(), "telegram query data invalid", "query", c.Request.URL.RawQuery)
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad query"})
		return
	}

	ctrl.log.DebugContext(c.Request.Context(), "telegram callback", "callback", callbackData)

	c.JSON(http.StatusOK, gin.H{})
}
