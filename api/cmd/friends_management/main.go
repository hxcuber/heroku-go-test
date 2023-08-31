package main

import (
	"context"
	"log"
	"os"
	"strings"

	router2 "github.com/hxcuber/friends-management/api/cmd/router"
	relationshipController "github.com/hxcuber/friends-management/api/internal/controller/relationship"
	systemController "github.com/hxcuber/friends-management/api/internal/controller/system"
	userController "github.com/hxcuber/friends-management/api/internal/controller/user"
	"github.com/hxcuber/friends-management/api/internal/repository"
	"github.com/hxcuber/friends-management/api/pkg/db/pg"
	"github.com/hxcuber/friends-management/api/pkg/env"
	"github.com/hxcuber/friends-management/api/pkg/httpserv"
)

const (
	host     = "localhost"
	port     = 5432
	username = "friends"
	password = "friends"
	dbname   = "friends"
)

func main() {
	ctx := context.Background()

	if err := run(ctx); err != nil {
		log.Fatal(err)
	}

	log.Println("Exiting...")
}

func run(ctx context.Context) error {
	log.Println("Starting app initialization")

	dbOpenConns := 4
	dbIdleConns := 2
	psqlInfo := env.GetAndValidateF("DB_URL")
	conn, err := pg.NewPool(psqlInfo, dbOpenConns, dbIdleConns)
	if err != nil {
		return err
	}

	defer conn.Close()

	rtr, _ := initRouter(ctx, conn)

	log.Println("App initialization completed")

	// err = httpserv.NewServer(rtr.Handler()).Start(ctx)
	// Using TempHandler for now to test flow, will refactor once I know
	// where the routing is meant to go
	err = httpserv.NewServer(rtr.Handler()).Start(ctx)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func initRouter(
	ctx context.Context,
	dbConn pg.BeginnerExecutor) (router2.Router, error) {
	registry := repository.New(dbConn)
	return router2.New(
		ctx,
		strings.Split(os.Getenv("CORS_ALLOWED_ORIGINS"), ","),
		os.Getenv("GQL_INTROSPECTION_ENABLED") == "true",
		systemController.New(registry),
		userController.New(registry),
		relationshipController.New(registry),
	), nil
}
