package db

import (
	"database/sql"
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2"
	_ "github.com/ClickHouse/clickhouse-go/v2"
	"go.uber.org/fx"
	"log/slog"
	"strconv"
	"time"
)

func getConnClickHouseStr(cfg *Config) string {
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

//func NewClickHouse(params ClickHouseParams) (*sql.DB, error) {
//	connStr := getConnClickHouseStr(&params.Config)
//	db, err := sql.Open("clickhouse", connStr)
//	if err != nil {
//		params.Logger.Error("error to open connection: " + err.Error())
//		return nil, fmt.Errorf("unable to connect to database: %w", err)
//	}
//
//	err = db.Ping()
//	if err != nil {
//		params.Logger.Error("error to ping database: " + err.Error())
//		return nil, fmt.Errorf("unable to ping to database: %w", err)
//	}
//
//	return db, nil
//}

func NewClickHouse(params ClickHouseParams) (*sql.DB, error) {
	conn := clickhouse.OpenDB(&clickhouse.Options{
		Addr: []string{params.Config.Host + ":" + strconv.Itoa(int(params.Config.Port))},
		Auth: clickhouse.Auth{
			Database: params.Config.DB,
			Username: params.Config.User,
			Password: params.Config.Password,
		},
		Settings: clickhouse.Settings{
			"max_execution_time": 60,
		},
		DialTimeout: time.Second * 30,
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		},
		Debug:                true,
		BlockBufferSize:      10,
		MaxCompressionBuffer: 10240,
		ClientInfo: clickhouse.ClientInfo{ // optional, please see Client info section in the README.md
			Products: []struct {
				Name    string
				Version string
			}{
				{Name: "my-app", Version: "0.1"},
			},
		},
	})
	conn.SetMaxIdleConns(5)
	conn.SetMaxOpenConns(10)
	conn.SetConnMaxLifetime(time.Hour)

	err := conn.Ping()
	if err != nil {
		return nil, err
	}

	return conn, nil
}
