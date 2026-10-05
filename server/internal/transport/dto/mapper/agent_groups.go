package mapper

import (
	"github.com/nougght/monitoring-system/server/internal/model/agent_groups"
	dto "github.com/nougght/monitoring-system/server/internal/transport/dto/types"
	"github.com/nougght/monitoring-system/shared/go/util"
)

func AgentGroupInfoToDTO(d *agent_groups.AgentGroupInfo) *dto.AgentGroupInfoDTO {
	if d == nil {
		return nil
	}
	return &dto.AgentGroupInfoDTO{
		Name:        d.Name,
		Description: d.Description,
	}
}

func AgentGroupInfoFromDTO(d *dto.AgentGroupInfoDTO) *agent_groups.AgentGroupInfo {
	if d == nil {
		return nil
	}
	return &agent_groups.AgentGroupInfo{
		Name:        d.Name,
		Description: d.Description,
	}
}

func CreateAgentGroupInputFromDTO(d *dto.CreateAgentGroupBody) *agent_groups.CreateAgentGroupInput {
	if d == nil {
		return nil
	}
	return &agent_groups.CreateAgentGroupInput{
		AgentGroupInfo: util.UnPtr(AgentGroupInfoFromDTO(&d.AgentGroupInfoDTO)),
		AgentIDs:       d.AgentIDs,
	}
}

func CreateAgentGroupResultToDTO(d *agent_groups.CreateAgentGroupResult) *dto.CreateAgentGroupResponse {
	if d == nil {
		return nil
	}
	return &dto.CreateAgentGroupResponse{
		AgentGroupDTO:    util.UnPtr(AgentGroupToDTO(&d.AgentGroup)),
		AgentsMovedCount: d.AgentsMovedCount,
	}
}

func UpdateAgentGroupInputFromDTO(d *dto.UpdateAgentGroupBody) *agent_groups.UpdateAgentGroupInput {
	if d == nil {
		return nil
	}
	return &agent_groups.UpdateAgentGroupInput{
		ID:             d.ID,
		AgentGroupInfo: util.UnPtr(AgentGroupInfoFromDTO(&d.AgentGroupInfoDTO)),
	}
}

func AgentGroupToDTO(d *agent_groups.AgentGroup) *dto.AgentGroupDTO {
	if d == nil {
		return nil
	}
	return &dto.AgentGroupDTO{
		ID:          d.ID,
		Name:        d.Name,
		Description: d.Description,
		CreatedAt:   d.CreatedAt,
		DeletedAt:   d.DeletedAt,
	}
}

func AgentGroupFromDTO(d *dto.AgentGroupDTO) *agent_groups.AgentGroup {
	if d == nil {
		return nil
	}
	return &agent_groups.AgentGroup{
		ID:          d.ID,
		Name:        d.Name,
		Description: d.Description,
		CreatedAt:   d.CreatedAt,
		DeletedAt:   d.DeletedAt,
	}
}
