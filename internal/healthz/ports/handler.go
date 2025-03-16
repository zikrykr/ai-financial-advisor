package ports

import "github.com/gin-gonic/gin"

type IHealthHandler interface {
	Healthz(c *gin.Context)
}
