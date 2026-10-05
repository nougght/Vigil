package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/nougght/monitoring-system/server/internal/model"
)

type AgentGroupInfoDTO struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
} // @Name AgentGroupInfo

type CreateAgentGroupBody struct {
	AgentGroupInfoDTO
	AgentIDs []uuid.UUID `json:"agentIDs"`
} // @Name CreateAgentGroupBody

type CreateAgentGroupResponse struct {
	AgentGroupDTO
	AgentsMovedCount int `json:"agentsMovedCount"`
} // @Name CreateAgentGroupResponse

type UpdateAgentGroupBody struct {
	ID          uuid.UUID              `json:"id"`
	Name        *string                `json:"name,omitempty"`
	Description model.Optional[string] `json:"description" swaggertype:"string"  extensions:"x-nullable"`
} // @Name UpdateAgentGroupBody

type AgentGroupDTO struct {
	ID          uuid.UUID  `json:"id"`
	Name        string     `json:"name"`
	Description *string    `json:"description,omitempty"`
	CreatedAt   time.Time  `json:"createdAt"`
	DeletedAt   *time.Time `json:"deletedAt,omitempty"`
} // @Name AgentGroup
