package rest

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nougght/monitoring-system/server/internal/service/overview"
	"github.com/nougght/monitoring-system/server/internal/transport/dto/mapper"
	_ "github.com/nougght/monitoring-system/server/internal/transport/dto/types"
)

type OverviewHandler struct {
	overviewService *overview.OverviewService
}

func NewOverviewHandler(overviewService *overview.OverviewService) *OverviewHandler {
	return &OverviewHandler{
		overviewService: overviewService,
	}
}

func (h *OverviewHandler) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/fleet/overview", h.GetFleetOverview)
}

// GetFleetOverview godoc
// @Id getFleetOverview
// @Summary Get fleet overview
// @Produce json
// @Success 200 {object} dto.AgentsOverview
// @Failure      400  {object}  map[string]any
// @Failure      404  {object}  map[string]any
// @Failure      500  {object}  map[string]any
// @Router /fleet/overview [get]
func (h *OverviewHandler) GetFleetOverview(c *gin.Context) {
	log.Println("overview")
	res, err := h.overviewService.GetOverview(c.Request.Context())
	if err != nil {
		handleError(c, fmt.Errorf("failed to get fleet overview: %w", err))
		return
	}

	c.JSON(http.StatusOK, mapper.OverviewToDTO(res))
}
