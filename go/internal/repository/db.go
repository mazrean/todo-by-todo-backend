package repository

//go:generate sqlc generate

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"time"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/mazrean/mazrean/todo-by-todo-backend/internal/repository/config"
	"github.com/mazrean/mazrean/todo-by-todo-backend/internal/repository/db"
)

type Repository struct {
	queries *db.Queries
	db      *sql.DB
}

type DBConfig struct {
	UserName string
	Password string
	Host     string
	Port     string
	Database string
}

func (c *DBConfig) DataSourceName(environment string) (string, error) {
	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return "", fmt.Errorf("failed to load location: %w", err)
	}

	tlsConfig := "false"
	if environment == "production" {
		tlsConfig = "true"
	}

	config := mysql.Config{
		User:                 c.UserName,
		Passwd:               c.Password,
		Net:                  "tcp",
		Addr:                 net.JoinHostPort(c.Host, c.Port),
		DBName:               c.Database,
		ParseTime:            true,
		Loc:                  loc,
		Collation:            "utf8mb4_general_ci",
		TLSConfig:            tlsConfig,
		AllowNativePasswords: true,
	}
	return config.FormatDSN(), nil
}

func NewDB(dbConf *DBConfig, env config.Environment) (*Repository, error) {
	dsn, err := dbConf.DataSourceName(string(env))
	if err != nil {
		return nil, fmt.Errorf("failed to get data source name: %w", err)
	}

	sqlDB, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	queries := db.New(sqlDB)

	return &Repository{
		queries: queries,
		db:      sqlDB,
	}, nil
}

func (r *Repository) Close() error {
	return r.db.Close()
}

func (r *Repository) GetQueries() *db.Queries {
	return r.queries
}

type contextKey string

const DBKey contextKey = "db"

func (r *Repository) Transaction(ctx context.Context, txOption *sql.TxOptions, fn func(ctx context.Context) error) error {
	tx, err := r.db.BeginTx(ctx, txOption)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	txCtx := context.WithValue(ctx, DBKey, tx)

	err = fn(txCtx)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *Repository) getDB(ctx context.Context) db.DBTX {
	if tx, ok := ctx.Value(DBKey).(*sql.Tx); ok {
		return tx
	}
	return r.db
}

func (r *Repository) GetQueriesWithTx(ctx context.Context) *db.Queries {
	return db.New(r.getDB(ctx))
}
