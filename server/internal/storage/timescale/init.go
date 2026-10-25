package timescale

import (
	"context"
	"embed"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nougght/monitoring-system/pkg/storage"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

func ConnectToDB(ctx context.Context, config *storage.PostgresConfig) (*pgxpool.Pool, error) {
	return storage.ConnectAndMigrate(ctx, config, migrationsFS)
}
