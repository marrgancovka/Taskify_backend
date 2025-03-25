package migrations

import (
	"database/sql"
	"embed"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/clickhouse"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"go.uber.org/fx"
	"log/slog"
)

// migrate -database "clickhouse://default:@localhost:9000/default" -path ./clickhouse/migrations up
// migrate create -ext sql -dir ./internal/pkg/migrations/schema -seq create_table_lists

type MigrationParams struct {
	fx.In

	DB     *sql.DB
	Logger *slog.Logger
}

//go:embed schema/*.sql
var migrationFiles embed.FS

func RunMigrations(p MigrationParams) error {
	sourceDriver, err := iofs.New(migrationFiles, "schema")
	if err != nil {
		p.Logger.Error("failed to load migration files: ", err.Error())
		return fmt.Errorf("failed to initialize migrations source driver: %w", err)
	}

	dbDriver, err := clickhouse.WithInstance(p.DB, &clickhouse.Config{})
	if err != nil {
		p.Logger.Error("failed to initialize clickhouse driver: ", err.Error())
		return fmt.Errorf("failed to initialize clickhouse driver: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", sourceDriver, "clickhouse", dbDriver)
	if err != nil {
		p.Logger.Error("failed to initialize migrate instance: ", err.Error())
		return fmt.Errorf("failed to initialize migrate instance: %w", err)
	}

	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		p.Logger.Error("failed to run migrations: ", err.Error())
		return fmt.Errorf("migration up failed: %w", err)
	}

	err = sourceDriver.Close()
	if err != nil {
		p.Logger.Error("failed to close migrations sourceDriver", "error", err)
		return fmt.Errorf("failed to close migrations sourceDriver: %w", err)
	}

	err = dbDriver.Close()
	if err != nil {
		p.Logger.Error("failed to close migrations dbDriver", "error", err)
		return fmt.Errorf("failed to close migrations dbDriver: %w", err)
	}

	return nil
}
