package agent_groups

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/nougght/monitoring-system/server/internal/config"
	"github.com/nougght/monitoring-system/server/internal/model"
	"github.com/nougght/monitoring-system/server/internal/model/agent_groups"
	agentregistry "github.com/nougght/monitoring-system/server/internal/service/agent_registry"
	"github.com/nougght/monitoring-system/server/internal/storage/timescale/repository"
)

type AgentGroupsService struct {
	cfg           *config.Config
	groupsRepo    *repository.AgentGroupsRepository
	agentRegistry *agentregistry.AgentRegistryService
	transactor    model.Transactor
}

func NewAgentGroupsService(cfg *config.Config, groupsRepo *repository.AgentGroupsRepository,
	agentRegistry *agentregistry.AgentRegistryService,
	transactor model.Transactor,
) (*AgentGroupsService, error) {
	if cfg == nil || groupsRepo == nil || transactor == nil {
		return nil, fmt.Errorf("params required")
	}
	return &AgentGroupsService{
		cfg:           cfg,
		groupsRepo:    groupsRepo,
		agentRegistry: agentRegistry,
		transactor:    transactor,
	}, nil
}

func (s *AgentGroupsService) CreateGroup(ctx context.Context, group *agent_groups.CreateAgentGroupInput) (*agent_groups.CreateAgentGroupResult, error) {
	var (
		createdGroup *agent_groups.AgentGroup
		moved        int
	)

	err := s.transactor.WithinTx(ctx, func(ctx context.Context) (err error) {
		createdGroup, err = s.groupsRepo.CreateGroup(ctx, &group.AgentGroupInfo)
		if errors.Is(err, repository.ErrConflict) {
			return model.ErrorAgentGroupNameIsTaken()
		}
		if err != nil {
			return fmt.Errorf("failed to create agent group: %w", err)
		}

		moved, err = s.MoveAgentsToGroup(ctx, group.AgentIDs, createdGroup.ID)
		if err != nil {
			return fmt.Errorf("failed to move agents to created group: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &agent_groups.CreateAgentGroupResult{
		AgentGroup:       *createdGroup,
		AgentsMovedCount: moved,
	}, nil
}

func (s *AgentGroupsService) GetGroupByID(ctx context.Context, id uuid.UUID) (*agent_groups.AgentGroup, error) {
	group, err := s.groupsRepo.GetGroupByID(ctx, id)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, model.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get agent group by ID: %w", err)
	}

	return group, nil
}

func (s *AgentGroupsService) GetAllGroups(ctx context.Context) ([]*agent_groups.AgentGroup, error) {
	groups, err := s.groupsRepo.GetAllGroups(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all agent groups: %w", err)
	}

	return groups, nil
}

func (s *AgentGroupsService) UpdateGroup(ctx context.Context, group *agent_groups.UpdateAgentGroupInput) error {
	err := s.groupsRepo.UpdateGroup(ctx, group)
	if errors.Is(err, repository.ErrConflict) {
		return model.ErrorAgentGroupNameIsTaken()
	}
	if errors.Is(err, repository.ErrNoAffectedRows) {
		return fmt.Errorf("not found group to update: %w", model.ErrNotFound)
	}
	if err != nil {
		return fmt.Errorf("failed to update agent group: %w", err)
	}

	return nil
}

func (s *AgentGroupsService) DeleteGroupByID(ctx context.Context, id uuid.UUID) error {
	err := s.groupsRepo.DeleteGroupByID(ctx, id)
	if errors.Is(err, repository.ErrNoAffectedRows) {
		return fmt.Errorf("not found group to delete: %w", model.ErrNotFound)
	}
	if err != nil {
		return fmt.Errorf("failed to delete agent group: %w", err)
	}

	return nil
}

func (s *AgentGroupsService) MoveAgentsToGroup(ctx context.Context, agentIDs []uuid.UUID, groupID uuid.UUID) (moved int, err error) {
	err = s.transactor.WithinTx(ctx, func(ctx context.Context) (err error) {
		group, err := s.groupsRepo.GetGroupByID(ctx, groupID)
		if errors.Is(err, repository.ErrNotFound) {
			return fmt.Errorf("not found group to move agents: %w", model.ErrNotFound)
		}
		if err != nil {
			return fmt.Errorf("failed to get group to move: %w", err)
		}

		moved, err = s.moveAgentsToGroup(ctx, agentIDs, group.ID)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return moved, nil
}

func (s *AgentGroupsService) MoveAgentToGroup(ctx context.Context, agentID uuid.UUID, groupID uuid.UUID) (err error) {
	err = s.transactor.WithinTx(ctx, func(ctx context.Context) (err error) {
		group, err := s.groupsRepo.GetGroupByID(ctx, groupID)
		if errors.Is(err, repository.ErrNotFound) {
			return fmt.Errorf("not found group to move agents: %w", model.ErrNotFound)
		}
		if err != nil {
			return fmt.Errorf("failed to get group to move: %w", err)
		}

		moved, err := s.moveAgentsToGroup(ctx, []uuid.UUID{agentID}, group.ID)
		if err != nil {
			return err
		}
		if moved == 0 {
			return fmt.Errorf("agent not found: %w", model.ErrNotFound)
		}

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *AgentGroupsService) moveAgentsToGroup(ctx context.Context, agentIDs []uuid.UUID, groupID uuid.UUID) (moved int, err error) {
	moved, err = s.agentRegistry.UpdateAgentsGroup(ctx, agentIDs, groupID)
	if err != nil {
		return 0, fmt.Errorf("failed to move agents to group: %w", err)
	}
	return moved, nil
}
