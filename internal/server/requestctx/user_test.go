package requestctx_test

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/ksusonic/kanban/internal/models"
	"github.com/ksusonic/kanban/internal/server/requestctx"
)

func TestSetUser_GetUser(t *testing.T) {
	c := &gin.Context{
		Request:  nil,
		Writer:   nil,
		Params:   nil,
		Keys:     nil,
		Errors:   nil,
		Accepted: nil,
	}

	t.Run("panic no user", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("panic assert")
			}
		}()

		requestctx.MustGetUser(c)
	})

	t.Run("success", func(t *testing.T) {
		identity := &models.UserIdentity{
			UserID: 123,
		}
		requestctx.SetUser(c, identity)

		assert.Equal(t, identity, requestctx.MustGetUser(c))
	})
}
