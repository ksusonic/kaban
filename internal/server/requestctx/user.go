package requestctx

import (
	"github.com/gin-gonic/gin"

	"github.com/ksusonic/kanban/internal/models"
)

func SetUser(c *gin.Context, identity *models.UserIdentity) {
	c.Set(userContextKey, identity)
}

func MustGetUser(c *gin.Context) *models.UserIdentity {
	value, ok := c.Get(userContextKey)
	if !ok {
		panic("no user in context")
	}

	if identity, castOK := value.(*models.UserIdentity); castOK {
		return identity
	}

	panic("no user in context")
}
