package main

import (
	"log"
	"nvrs-gateway/config"
	"nvrs-gateway/routes"
	"nvrs-gateway/storage"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	config.Load() // Load environment-specific configuration

	// Initialize the database based on the environment
	if config.Env == "production" {
		if err := storage.InitSQLite(); err != nil {
			log.Fatalf("Failed to initialize SQLite: %v", err)
		}
		// if err := storage.InitPostgres(config.DB_URL); err != nil {
		//     log.Fatalf("Failed to initialize PostgreSQL: %v", err)
		// }
	} else {
		if err := storage.InitSQLite(); err != nil {
			log.Fatalf("Failed to initialize SQLite: %v", err)
		}
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	routes.InitRoutes(e)

	e.Logger.Fatal(e.Start(":8080"))
}
