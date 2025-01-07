package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const tgCallbackPath = "/auth/tg-callback"

func (ctrl *Controller) Page(c *gin.Context) {
	c.Header("Content-Security-Policy", "default-src *; style-src 'self' 'unsafe-inline'; script-src 'self' 'unsafe-inline' 'unsafe-eval' https://telegram.org; frame-ancestors 'self' https://oauth.telegram.org;") //nolint:lll
	c.HTML(http.StatusOK, "login.tmpl", gin.H{
		"next":          c.DefaultQuery("next", "/"),
		"bot_name":      ctrl.botCfg.Name,
		"callback_path": tgCallbackPath,
	})
}
