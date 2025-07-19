package main

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/mazrean/mazrean/todo-by-todo-backend/internal/di"
	"github.com/mazrean/mazrean/todo-by-todo-backend/internal/router"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

type CLI struct {
	LogLevel string           `kong:"short='l',help='Log level',enum='debug,info,warn,error',default='info'"`
	Version  kong.VersionFlag `kong:"short='v',help='Show version and exit.'"`
	Addr     string           `kong:"help='Address to run the server on',default=':8080',env='ADDR'"`
}

func (cli *CLI) Run() error {
	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: parseLogLevel(cli.LogLevel),
	})
	logger := slog.New(handler)
	slog.SetDefault(logger)

	app := di.InjectCLI(
		router.Addr(cli.Addr),
		router.Version(version),
	)

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
