package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	logger "go.uber.org/zap"

	ahttp "webChat/internal/api/http/handler"
	"webChat/internal/config"
	"webChat/internal/db"
	"webChat/internal/db/store"

	"github.com/oklog/run"
)

func main() {

	var (
		cfgPath string
		ctx     = context.Background()
		log, _  = logger.NewDevelopment()
		logger  = log.Sugar()
	)
	defer logger.Sync()

	cfgPath, ok := os.LookupEnv("CONFIG_PATH")
	if !ok {
		cfgPath = "./config.json"
	}

	cfg, err := config.NewConfig(cfgPath)
	if err != nil {
		logger.Fatal("failed to create a new config", err)
	}

	dbConn, err := db.NewConnection(&cfg.DB)
	if err != nil {
		logger.Fatal("failed to connect to db", err)

	}
	defer dbConn.Close()

	logger.Debug("database successfully connected")

	userStore := store.NewUserStore(dbConn, "users")

	var g run.Group

	ahttp.NewServer(cfg, logger, userStore).Run(&g)

	ctx, cancel := context.WithCancel(context.Background())
	g.Add(func() error {
		shutdown := make(chan os.Signal, 1)
		signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

		logger.Info("[signal-watcher] started")

		select {
		case sig := <-shutdown:
			return fmt.Errorf("terminated with signal: %s", sig.String())
		case <-ctx.Done():
			return nil
		}
	}, func(err error) {
		cancel()
		logger.Error("gracefully shutdown application", err)
	})

	logger.Error("application stopped", g.Run())
}
