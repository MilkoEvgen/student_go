package test

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/docker/go-connections/nat"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	pgContainer testcontainers.Container
	ctx         context.Context
	cancel      context.CancelFunc
)

func StartTestContainer() {
	ctx, cancel = context.WithCancel(context.Background())

	pg, err := postgres.RunContainer(
		ctx,
		testcontainers.WithImage("postgres:16"),
		postgres.WithDatabase("test"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForSQL("5432/tcp", "postgres", func(host string, port nat.Port) string {
				return fmt.Sprintf("host=%s port=%s user=test password=test dbname=test sslmode=disable", host, port.Port())
			}).WithStartupTimeout(30*time.Second),
		),
	)
	if err != nil {
		log.Fatalf("[test] failed to start container: %v", err)
	}

	pgContainer = pg

	host, _ := pg.Host(ctx)
	port, _ := pg.MappedPort(ctx, "5432")

	os.Setenv("APP_ENV", "test")
	os.Setenv("DB_HOST", host)
	os.Setenv("DB_PORT", port.Port())
	os.Setenv("DB_NAME", "test")
	os.Setenv("DB_USER", "test")
	os.Setenv("DB_PASSWORD", "test")

	log.Println("[test] started postgres container")
}

func StopTestContainer() {
	if pgContainer != nil {
		if err := pgContainer.Terminate(ctx); err != nil {
			log.Printf("[test] failed to stop container: %v", err)
		} else {
			log.Println("[test] stopped postgres container")
		}
	}
	if cancel != nil {
		cancel()
	}
}
