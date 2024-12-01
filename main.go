package main

import (
	"database/sql"
	"loan-management/propose"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	cfg, err := NewConfigFromEnv()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	db, err := NewDB(cfg.ConnectionString)
	if err != nil {
		log.Fatalf("error connecting to db: %v", err)
	}

	start(cfg, db)
}

func start(cfg Config, db *sql.DB) {
	e := echo.New()
	e.GET("/healthz", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "ok",
			"port":   cfg.Port,
		})
	})

	repo := propose.NewRepository(db)
	command := propose.NewCommand(repo)
	proposeHandler := propose.NewHandler(command)
	e.POST("/propose", proposeHandler.Handle)

	e.Logger.Fatal(e.Start(":" + cfg.Port))
}

func NewDB(connectionString string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}
	return db, nil
}
