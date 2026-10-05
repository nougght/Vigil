package util

import (
	"context"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nougght/monitoring-system/server/internal/model"
)

func GetAgentIDFromContext(ctx context.Context) (uuid.UUID, error) {
	if ctx == nil {
		return uuid.Nil, fmt.Errorf("ctx is nil")
	}
	idFromTLS := ctx.Value(model.ContextKeyAgentID)
	if idFromTLS == nil {
		return uuid.Nil, nil
	}
	id, ok := idFromTLS.(uuid.UUID)
	if !ok {
		return uuid.Nil, errors.New("invalid agent ID in context")
	}
	return id, nil
}

func GetParamID(c *gin.Context, name string) (uuid.UUID, error) {
	return uuid.Parse(c.Param(name))
}
