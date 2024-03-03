package main

import (
	_ "github.com/lib/pq"
	"log"
	"os"
	psql_helper "postgresHelper/internal/app/psql-helper"
	"postgresHelper/internal/collector"
	"postgresHelper/internal/config"
	"postgresHelper/internal/runner"
	"postgresHelper/internal/storage"
)

func main() {
	defer func() {
		recover()
	}()

	err := config.Init()
	if err != nil {
		log.Fatal(err)
	}

	conn, err := createPostgresConn(&config.ConfigStruct.PG)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	store := storage.New()
	run := runner.New(collector.NewCollector(conn), store)
	run.Run()

	delivery := psql_helper.New(store)
	grpcServer, err := runGRPCServer(delivery, &config.ConfigStruct.GRPC)
	if err != nil {
		log.Fatal(err)
	}
	defer grpcServer.Close()
	Lock(make(chan os.Signal, 1))
}
