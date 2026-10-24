package agent_groups

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/nougght/monitoring-system/server/internal/model"
)

type (
	AgentGroup struct {
		ID          uuid.UUID  `db:"id"`
		Name        string     `db:"name"`
		Description *string    `db:"description"`
		CreatedAt   time.Time  `db:"created_at"`
		UpdatedAt   time.Time  `db:"updated_at"`
		DeletedAt   *time.Time `db:"deleted_at"`
	}

	AgentGroupInfo struct {
		Name        string  `db:"name"`
		Description *string `db:"description"`
	}

	CreateAgentGroupInput struct {
		AgentGroupInfo
		AgentIDs []uuid.UUID `db:"agent_ids"`
	}

	CreateAgentGroupResult struct {
		AgentGroup
		AgentsMovedCount int
	}

	UpdateAgentGroupInput struct {
		ID          uuid.UUID
		Name        *string                `db:"name"`
		Description model.Optional[string] `db:"description"`
	}
)

func (g *CreateAgentGroupInput) Validate() error {
	return validateName(g.Name)
}

func (g *UpdateAgentGroupInput) Validate() error {
	if g.Name != nil {
		return validateName(*g.Name)
	}
	return nil
}

func validateName(name string) error {
	if strings.Contains(name, " ") {
		return fmt.Errorf("name must not contain spaces: %w", model.ErrBadRequest)
	}
	if name == "" {
		return fmt.Errorf("name must not be empty: %w", model.ErrBadRequest)
	}
	return nil
}
