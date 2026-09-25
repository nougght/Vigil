package agent_groups

import (
	"time"

	"github.com/google/uuid"
	"github.com/nougght/monitoring-system/server/internal/model"
)

type (
	AgentGroup struct {
		ID          uuid.UUID `db:"id"`
		Name        string    `db:"name"`
		Description *string   `db:"description"`
		CreatedAt   time.Time `db:"created_at"`
		UpdatedAt   time.Time `db:"updated_at"`
	}

	AgentGroupInput struct {
		Name        string                 `db:"name"`
		Description model.Optional[string] `db:"description"`
	}

	CreateAgentGroupInput struct {
		AgentGroupInput
		AgentIDs []uuid.UUID `db:"agent_ids"`
	}
	UpdateAgentGroupInput struct {
		ID uuid.UUID
		AgentGroupInput
	}
)
