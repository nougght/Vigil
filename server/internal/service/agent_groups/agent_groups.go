package agentregistry

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/jackc/pgx/v4"
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
	transactor model.Transactor, cert *model.Certs,
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
	tx, err := s.transactor.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed begin transaction: %w", err)
	}
	defer func() {
		err := tx.Rollback(ctx)
		if err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			log.Printf("rollback failed: %s", err.Error())
		}
	}()
	ctx = context.WithValue(ctx, model.ContextKeyTx, tx)

	createdGroup, err := s.groupsRepo.CreateGroup(ctx, &group.AgentGroupInput)
	if errors.Is(err, repository.ErrConflict) {
		return nil, model.ErrorAgentGroupNameIsTaken()
	}
	if err != nil {
		return nil, fmt.Errorf("failed to create agent group: %w", err)
	}

	moved, err := s.agentRegistry.MoveAgentsToGroup(ctx, group.AgentIDs, createdGroup.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to move agents to created group: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, fmt.Errorf("commit transaction error: %w", err)
	}

	return &agent_groups.CreateAgentGroupResult{
		AgentGroup:       *createdGroup,
		AgentsMovedCount: moved,
	}, nil
}
