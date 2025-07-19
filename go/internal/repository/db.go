package repository

//go:generate sqlc generate

import (
	"database/sql"
	"fmt"
	"net"
	"time"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
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

func (c *DBConfig) DataSourceName() (string, error) {
	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return "", fmt.Errorf("failed to load location: %w", err)
	}

	config := mysql.Config{
		User:      c.UserName,
		Passwd:    c.Password,
		Net:       "tcp",
		Addr:      net.JoinHostPort(c.Host, c.Port),
		DBName:    c.Database,
		ParseTime: true,
		Loc:       loc,
		Collation: "utf8mb4_general_ci",
	}
	return config.FormatDSN(), nil
}

func NewDB(dbConf *DBConfig) (*Repository, error) {
	dsn, err := dbConf.DataSourceName()
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
