package service

import (
	"log"

	"github.com/nougght/monitoring-system/server/internal/config"
	"github.com/nougght/monitoring-system/server/internal/infrastructure/eventbus"
	"github.com/nougght/monitoring-system/server/internal/model"
	agent "github.com/nougght/monitoring-system/server/internal/service/agent_interaction"
	agentregistry "github.com/nougght/monitoring-system/server/internal/service/agent_registry"
	"github.com/nougght/monitoring-system/server/internal/service/metrics"
	"github.com/nougght/monitoring-system/server/internal/service/realtime"
	"github.com/nougght/monitoring-system/server/internal/storage/timescale/repository"
)

type Services struct {
	agentRegistry *agentregistry.AgentRegistryService
	agent         *agent.AgentInteractionService
	metrics       *metrics.MetricsService
	realtime      *realtime.RealtimeService
}

//nolint:unused
type ServicesOptions struct {
	Config       *config.Config
	Repositories *repository.Repositories
	Transactor   model.Transactor
	Cert         *model.Certs
	Bus          *eventbus.EventBus
}

func New(opts ServicesOptions) *Services {
	agentRegistry, err := agentregistry.NewAgentRegistryService(
		opts.Config,
		opts.Repositories.AgentRepository(),
		opts.Repositories.EnrollmentKeysRepository(),
		opts.Repositories.SpecsRepository(),
		opts.Transactor,
		opts.Cert,
	)
	if err != nil {
		log.Panicf("failed initialize agent registry: %s", err.Error())
	}
	metrics, err := metrics.NewMetricsService(
		opts.Config,
		opts.Transactor,
		opts.Repositories.MetricsRepository(),
		opts.Repositories.SeriesRepository(),
		opts.Bus,
	)
	agentInteraction, err := agent.NewAgentInteractionService(
		opts.Config,
		agentRegistry,
		metrics,
	)
	if err != nil {
		log.Panicf("failed initialize agent interaction service: %s", err.Error())
	}
	realtimeService, err := realtime.NewRealtimeService(
		opts.Config,
		opts.Transactor,
		opts.Bus,
	)
	return &Services{
		agentRegistry: agentRegistry,
		agent:         agentInteraction,
		metrics:       metrics,
		realtime:      realtimeService,
	}
}

func (s *Services) AgentRegistry() *agentregistry.AgentRegistryService {
	return s.agentRegistry
}

func (s *Services) AgentInteractionService() *agent.AgentInteractionService {
	return s.agent
}

func (s *Services) Metrics() *metrics.MetricsService {
	return s.metrics
}

func (s *Services) Realtime() *realtime.RealtimeService {
	return s.realtime
}
