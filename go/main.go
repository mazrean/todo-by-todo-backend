package main

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/mazrean/mazrean/todo-by-todo-backend/internal/di"
	"github.com/mazrean/mazrean/todo-by-todo-backend/internal/repository"
	"github.com/mazrean/mazrean/todo-by-todo-backend/internal/repository/config"
	"github.com/mazrean/mazrean/todo-by-todo-backend/internal/router"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

type CLI struct {
	LogLevel    string           `kong:"short='l',help='Log level',enum='debug,info,warn,error',default='info'"`
	Version     kong.VersionFlag `kong:"short='v',help='Show version and exit.'"`
	Addr        string           `kong:"help='Address to run the server on',default=':8080',env='ADDR'"`
	Environment string           `kong:"help='Environment',enum='development,production',default='development',env='ENV'"`
	DBUserName  string           `kong:"help='Database user name',default='root',env='DB_USER_NAME'"`
	DBPassword  string           `kong:"help='Database password',default='password',env='DB_PASSWORD'"`
	DBHost      string           `kong:"help='Database host',default='localhost',env='DB_HOST'"`
	DBPort      string           `kong:"help='Database port',default='3306',env='DB_PORT'"`
	DBDatabase  string           `kong:"help='Database name',default='todo',env='DB_DATABASE'"`
}

func (cli *CLI) Run() error {
	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: parseLogLevel(cli.LogLevel),
	})
	logger := slog.New(handler)
	slog.SetDefault(logger)

	dbCfg := &repository.DBConfig{
		UserName: cli.DBUserName,
		Password: cli.DBPassword,
		Host:     cli.DBHost,
		Port:     cli.DBPort,
		Database: cli.DBDatabase,
	}

	environment := config.Environment(cli.Environment)

	app, err := di.InjectCLI(
		router.Addr(cli.Addr),
		router.Version(version),
		dbCfg,
		environment,
	)

	if err != nil {
		return fmt.Errorf("failed to inject dependencies: %w", err)
	}

	return app.Run()
}

func main() {
	var cli CLI
	kongCtx := kong.Parse(&cli,
		kong.Name("app"),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
		}),
		kong.Vars{
			"version": fmt.Sprintf("%s (%s) released on %s", version, commit, date),
		},
	)

	if err := kongCtx.Run(); err != nil {
		panic(err)
	}
}

func parseLogLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
