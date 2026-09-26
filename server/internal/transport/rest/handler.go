package rest

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nougght/monitoring-system/server/internal/model"
	"github.com/nougght/monitoring-system/server/internal/service"
	dto "github.com/nougght/monitoring-system/server/internal/transport/dto/types"
)

// write http response based on service error
func handleError(c *gin.Context, err error) {
	serr := &model.ServiceError{}
	if !errors.As(err, &serr) {
		serr = model.NewError(err)
	}
	responseCode := http.StatusInternalServerError
	resp := dto.ErrorResponse{
		Key:     serr.Key,
		Message: "internal server error",
	}

	switch {
	case errors.Is(serr.Err, model.ErrUnauthorized):
		responseCode = http.StatusUnauthorized
	case errors.Is(serr.Err, model.ErrBadRequest):
		responseCode = http.StatusBadRequest
	case errors.Is(serr.Err, model.ErrNotFound):
		responseCode = http.StatusNotFound
	case errors.Is(serr.Err, model.ErrServiceUnavailable):
		responseCode = http.StatusServiceUnavailable
	case errors.Is(err, model.ErrServiceUnavailable):
		responseCode = http.StatusServiceUnavailable
	}

	if responseCode != http.StatusInternalServerError {
		resp.Message = err.Error()
	}

	log.Printf("Error: %v", err)
	c.JSON(responseCode, resp)
}

type Handlers struct {
	agentHandler    *AgentHandler
	streamHandler   *StreamHandler
	overviewHandler *OverviewHandler
}

func newHandlers(services *service.Services) *Handlers {
	agent := newAgentHandler(
		services.AgentRegistry(),
		services.AgentInteractionService(),
	)

	stream := newStreamHandler(
		services.AgentRegistry(),
		services.AgentInteractionService(),
	)
	overview := NewOverviewHandler(
		services.Overview(),
	)

	return &Handlers{
		agentHandler:    agent,
		streamHandler:   stream,
		overviewHandler: overview,
	}
}

func (h *Handlers) AgentHandler() *AgentHandler {
	return h.agentHandler
}

func (h *Handlers) StreamHandler() *StreamHandler {
	return h.streamHandler
}

func (h *Handlers) OverviewHandler() *OverviewHandler {
	return h.overviewHandler
}
