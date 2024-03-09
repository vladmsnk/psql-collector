package main

import (
	"context"
	"fmt"
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
	collect := collector.NewCollector(conn)
	run := runner.New(collect, store)
	run.Run()

	delivery := psql_helper.New(store)
	grpcServer, err := runGRPCServer(delivery, &config.ConfigStruct.GRPC)
	if err != nil {
		log.Fatal(err)
	}
	defer grpcServer.Close()

	res, err := collect.CollectQueryTypesDistribution(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(res)
	//loader := loader.New(conn)
	//loader.FillWithTestData(context.Background())

	Lock(make(chan os.Signal, 1))
}
