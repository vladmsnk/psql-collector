package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"postgresHelper/internal/collector"
	"postgresHelper/internal/runner"
	"postgresHelper/internal/storage"
	"time"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "user"
	password = "password"
	dbname   = "postgres"
)

func createPostgresConn() (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
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

func main() {
	conn, err := createPostgresConn()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	cl := collector.NewCollector(conn)
	st := storage.New()

	run := runner.New(cl, st)
	run.Run(context.Background())

	time.Sleep(time.Minute * 10)
}
