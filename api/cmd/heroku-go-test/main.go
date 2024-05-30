package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/hxcuber/friends-management/api/cmd/router"
	relationshipController "github.com/hxcuber/friends-management/api/internal/controller/relationship"
	systemController "github.com/hxcuber/friends-management/api/internal/controller/system"
	userController "github.com/hxcuber/friends-management/api/internal/controller/user"
	"github.com/hxcuber/friends-management/api/internal/repository"
	"github.com/hxcuber/friends-management/api/pkg/app"
	"github.com/hxcuber/friends-management/api/pkg/db/pg"
	"github.com/hxcuber/friends-management/api/pkg/env"
	"github.com/hxcuber/friends-management/api/pkg/httpserv"
	"github.com/pkg/errors"
	"strconv"
)

func main() {
	ctx := context.Background()

	appCfg := app.Config{
		ProjectName:      env.GetAndValidateF("PROJECT_NAME"),
		AppName:          env.GetAndValidateF("APP_NAME"),
		SubComponentName: env.GetAndValidateF("PROJECT_COMPONENT"),
		Env:              app.Env(env.GetAndValidateF("APP_ENV")),
		Version:          env.GetAndValidateF("APP_VERSION"),
		Server:           env.GetAndValidateF("SERVER_NAME"),
		ProjectTeam:      os.Getenv("PROJECT_TEAM"),
	}
	if err := appCfg.IsValid(); err != nil {
		log.Fatal(err)
	}

	if err := run(ctx); err != nil {
		log.Fatal(err)
	}

	log.Println("Exiting...")
}

func run(ctx context.Context) error {
	log.Println("Starting app initialization")

	dbOpenConns, err := strconv.Atoi(env.GetAndValidateF("DB_POOL_MAX_OPEN_CONNS"))
	if err != nil {
		return errors.WithStack(fmt.Errorf("invalid db pool max open conns: %w", err))
	}
	dbIdleConns, err := strconv.Atoi(env.GetAndValidateF("DB_POOL_MAX_IDLE_CONNS"))
	if err != nil {
		return errors.WithStack(fmt.Errorf("invalid db pool max idle conns: %w", err))
	}

	conn, err := pg.NewPool(env.GetAndValidateF("DATABASE_URL"), dbOpenConns, dbIdleConns)

	if err != nil {
		return err
	}

	defer conn.Close()

	rtr, _ := initRouter(ctx, conn)

	log.Println("App initialization completed")

	httpserv.NewServer(rtr.Handler()).Start(ctx)

	return nil
}

func initRouter(
	ctx context.Context,
	dbConn pg.BeginnerExecutor) (router.Router, error) {
	registry := repository.New(dbConn)
	return router.New(
		ctx,
		strings.Split(os.Getenv("CORS_ALLOWED_ORIGINS"), ","),
		os.Getenv("GQL_INTROSPECTION_ENABLED") == "true",
		systemController.New(registry),
		userController.New(registry),
		relationshipController.New(registry),
	), nil
}
