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

	assert.Nil(t, requestctx.GetUser(c))

	identity := &models.UserIdentity{UserID: 123}
	requestctx.SetUser(c, identity)

	assert.Equal(t, identity, requestctx.GetUser(c))
}
