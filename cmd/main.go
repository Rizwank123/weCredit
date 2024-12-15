package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/weCredit/internal/dependency"
	"github.com/weCredit/internal/http/swagger"
	"github.com/weCredit/internal/pkg/config"
)

func main() {
	// Create a new map
	cfgOpt := getConfigOptions()
	cfg, err := dependency.NewConfig(cfgOpt)
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	// setup database connections
	// Setup database connection
	db, err := dependency.NewDatabaseConfig(cfg)
	if err != nil {
		log.Fatalf("failed to create connection for database: %v", err)
	}
	err = db.Ping(context.Background())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()
	// initialize the dependencies
	api, err := dependency.NewWeCredit(cfg, db)
	if err != nil {
		log.Fatalf("failed to create weCredit instance: %v", err)
	}
	// setup echo framework
	e := echo.New()
	//setup middleware
	api.SetupMiddleware(e)
	//setup swagger
	swagger.SetupSwagger(cfg, e)
	//setup routes
	api.SetupRoutes(e)
	//setup server in a goroutine
	go func() {
		e.Logger.Info(e.Start(fmt.Sprintf("0.0.0.0:%d", cfg.AppPort)))
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Server is shuting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
	log.Println("Server successFully shut down")

}

func getConfigOptions() config.Options {
	cfgSource := os.Getenv(config.SourceKey)
	if cfgSource == "" {
		cfgSource = config.SourceEnv
	}
	return config.Options{
		ConfigFileSource: cfgSource,
		ConfigFile:       ".env",
	}
}
