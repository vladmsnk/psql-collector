package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"postgresHelper/internal/config"
	desc "postgresHelper/internal/pkg/collector"
	"postgresHelper/lib/grpc_server"
	"strconv"
	"syscall"
)

func runGRPCServer(implementation desc.CollectorServer, cfg *grpc_server.GRPCConfig) (*grpc_server.GRPCServer, error) {
	grpcServer, err := grpc_server.NewGRPCServer(cfg)
	if err != nil {
		return nil, fmt.Errorf("grpc_server.NewGRPCServer: %w", err)
	}

	desc.RegisterCollectorServer(grpcServer.Ser, implementation)
	grpcServer.Run()

	log.Printf("started grpc server at %s:%s", cfg.Host, strconv.Itoa(cfg.Port))
	return grpcServer, nil
}

func createPostgresConn(cfg *config.Postgres) (*sql.DB, error) {
	psqlInfo := cfg.GetConnectionInfo()
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("error pinging database: %v", err)
	}

	log.Println("Successfully connected to postgres database")
	return db, nil
}

func Lock(ch chan os.Signal) {
	defer func() {
		ch <- os.Interrupt
	}()
	signal.Notify(ch,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	<-ch
}
