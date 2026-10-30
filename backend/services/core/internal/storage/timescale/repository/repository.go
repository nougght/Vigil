package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/nougght/monitoring-system/backend/pkg/storage"
)

type BaseRepository struct {
	db storage.DB
}

func (r *BaseRepository) conn(ctx context.Context) storage.DB {
	return storage.Conn(ctx, r.db)
}

type Repositories struct {
	agentRpository           *AgentRepository
	enrollmentKeysRepository *EnrollmentKeysRepository
	specsRepository          *SpecsRepository
	metricsRepository        *MetricsRepository
	seriesRepository         *SeriesRepository
	agentGroupsRepository    *AgentGroupsRepository
}

func New(db storage.DB) *Repositories {
	return &Repositories{
		agentRpository:           NewAgentRepository(db),
		enrollmentKeysRepository: NewEnrollmentKeysRepository(db),
		specsRepository:          NewSpecsRepository(db),
		metricsRepository:        NewMetricsRepository(db),
		seriesRepository:         NewSeriesRepository(db),
		agentGroupsRepository:    NewAgentGroupsRepository(db),
	}
}

func (r *Repositories) AgentRepository() *AgentRepository {
	return r.agentRpository
}

func (r *Repositories) EnrollmentKeysRepository() *EnrollmentKeysRepository {
	return r.enrollmentKeysRepository
}

func (r *Repositories) SpecsRepository() *SpecsRepository {
	return r.specsRepository
}

func (r *Repositories) MetricsRepository() *MetricsRepository {
	return r.metricsRepository
}

func (r *Repositories) SeriesRepository() *SeriesRepository {
	return r.seriesRepository
}
func (r *Repositories) AgentGroupsRepository() *AgentGroupsRepository {
	return r.agentGroupsRepository
}

func CollectOnePtr[T any](rows pgx.Rows) (*T, error) {
	return pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[T])
}

func CollectRowsPtr[T any](rows pgx.Rows) ([]*T, error) {
	return pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[T])
}
