package requestctx

import (
	"github.com/gin-gonic/gin"

	"github.com/ksusonic/kanban/internal/models"
)

func SetUser(c *gin.Context, identity *models.UserIdentity) {
	c.Set(userContextKey, identity)
}

func GetUser(c *gin.Context) *models.UserIdentity {
	if value, has := c.Get(userContextKey); has {
		if typedValue, ok := value.(*models.UserIdentity); ok {
			return typedValue
		}
	}

	return nil
}
