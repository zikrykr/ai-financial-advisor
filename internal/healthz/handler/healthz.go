package handler

import (
	"net/http"

	"github.com/ai-financial-advisor/internal/healthz/response"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"gorm.io/gorm"
)

type HealthHandler struct {
	db *gorm.DB
}

func NewHealthHandler(db *gorm.DB) HealthHandler {
	return HealthHandler{db: db}
}

// Healthz godoc
// @Summary Check health of the service
// @Description Check health of the service
// @Tags healthz
// @Accept json
// @Produce json
// @Success 200 {object} string
// @Router /healthz [get]
func (h HealthHandler) Healthz(c *gin.Context) {
	tracer := otel.Tracer("health-handler")
	ctx, span := tracer.Start(c.Request.Context(), "Healthz")
	defer span.End()

	res := response.HealthResponse{}

	db, err := h.db.DB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.HealthResponse{
			Status: "unhealthy",
			DB:     "disconnected",
		})
		return
	}

	if err := db.PingContext(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, response.HealthResponse{
			Status: "unhealthy",
			DB:     "disconnected",
		})
		return
	} else {
		res.DB = "connected"
		res.Status = "healthy"
	}

	c.JSON(http.StatusOK, res)
}
