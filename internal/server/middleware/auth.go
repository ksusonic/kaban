package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ksusonic/kanban/internal/models"
	"github.com/ksusonic/kanban/internal/server/api"
	"github.com/ksusonic/kanban/internal/server/requestctx"
)

type AuthProvider interface {
	CheckToken(tokenString string) (*models.UserIdentity, error)
}

func AuthRequired(authProvider AuthProvider) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenValue := c.GetHeader("Authorization")
		userIdentity, err := authProvider.CheckToken(tokenValue)
		if err != nil {
			c.JSON(http.StatusUnauthorized, api.ErrorResponse{
				Error: "authorization required"},
			)
			c.Abort()
			return
		}

		requestctx.SetUser(c, userIdentity)

		c.Next()
	}
}
