package timescale

import (
	"context"
	"embed"

	"github.com/nougght/monitoring-system/backend/pkg/storage"

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

func ConnectToDB(ctx context.Context, config *storage.PostgresConfig) (*pgxpool.Pool, error) {
	return storage.ConnectAndMigrate(ctx, config, migrationsFS)
}
