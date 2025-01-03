package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (ctrl *Controller) Login(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}
