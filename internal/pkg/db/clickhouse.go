package db

import (
	"database/sql"
	"fmt"
	_ "github.com/ClickHouse/clickhouse-go/v2"
	"go.uber.org/fx"
	"log/slog"
)

func getConnStr(cfg *Config) string {
	return fmt.Sprintf(
		"clickhouse://%s:%s@%s:%d/%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DB,
	)
}

type ClickHouseParams struct {
	fx.In

	Logger *slog.Logger
	Config Config
}

func NewClickHouse(params ClickHouseParams) (*sql.DB, error) {
	connStr := getConnStr(&params.Config)
	db, err := sql.Open("clickhouse", connStr)
	if err != nil {
		params.Logger.Error("error to open connection: " + err.Error())
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	err = db.Ping()
	if err != nil {
		params.Logger.Error("error to ping database: " + err.Error())
		return nil, fmt.Errorf("unable to ping to database: %w", err)
	}

	return db, nil
}
