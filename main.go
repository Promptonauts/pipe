package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Promptonauts/pipe/pkg/api"
	"github.com/Promptonauts/pipe/pkg/controlplane"
	"github.com/Promptonauts/pipe/pkg/executor"
	"github.com/Promptonauts/pipe/pkg/guardrails"
	"github.com/Promptonauts/pipe/pkg/observability"
	"github.com/Promptonauts/pipe/pkg/scheduler"
	"github.com/Promptonauts/pipe/pkg/store"
)

func main() {
	logger := observability.NewLogger("pipe-server")
	metrics := observability.NewMetricsRegistry()

	db, err := store.NewSQLiteStore("pipe.db")
	if err != nil {
		log.Fatalf("failed to initialize store: %v", err)
	}
	defer db.Close()

	if err := db.Migrate(); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	guardrailEngine := guardrails.NewEngine(metrics, logger)
	execEngine := executor.NewEngine(db, guardrailEngine, metrics, logger)
	sched := scheduler.NewScheduler(execEngine, logger, metrics, scheduler.Config{
		MaxConcurrency:    10,
		QueueSize:         1000,
		WorkerCount:       5,
		OverloadThreshold: 800,
	})

	reconciler := controlplane.NewReconciler(db, execEngine, sched, logger, metrics)
	controller := controlplane.NewController(reconciler, db, logger)

	go controller.Run()
	go sched.Start()

	srv := api.NewServer(db, controller, sched, metrics, logger)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigCh
		logger.Info("shutting down...")
		sched.Stop()
		controller.Stop()
		os.Exit(0)
	}()

	logger.Info("PIPE server starting on :8080")
	if err := srv.Run(":8080"); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
