package main

import (
	"database/sql"
	"loan-management/application/approve"
	"loan-management/application/invest"
	loan_query "loan-management/application/loan-query"
	"loan-management/application/propose"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	cfg, err := NewConfigFromEnv()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	db, err := NewDB()
	if err != nil {
		log.Fatalf("error connecting to db: %v", err)
	}

	start(cfg, db)
}

func start(cfg Config, db *sql.DB) {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.RequestID())

	e.GET("/healthz", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "ok",
			"port":   cfg.Port,
		})
	})

	registerPropose(db, e)
	registerApprove(db, e)
	registerLoanQuery(db, e)
	registerInvest(db, e)

	e.Logger.Fatal(e.Start(":" + cfg.Port))
}

func registerPropose(db *sql.DB, e *echo.Echo) {
	repo := propose.NewRepository(db)
	command := propose.NewCommand(repo)
	proposeHandler := propose.NewHandler(command)
	e.POST("/", proposeHandler.Handle)
}

func registerLoanQuery(db *sql.DB, e *echo.Echo) {
	queryHandler := loan_query.NewHandler(db)
	e.GET("/:loan_id", queryHandler.Handle)
}

func registerApprove(db *sql.DB, e *echo.Echo) {
	repo := approve.NewRepository(db)
	command := approve.NewCommand(repo)
	approveHandler := approve.NewHandler(command)
	e.POST("/approve/:loan_id", approveHandler.Handle)
}

func registerInvest(db *sql.DB, e *echo.Echo) {
	repo := invest.NewRepository(db)
	command := invest.NewCommand(repo)
	investHandler := invest.NewHandler(command)
	e.POST("/invest/:loan_id", investHandler.Handle)
}

func NewDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}

	db.Exec("CREATE TABLE loans (id UUID, borrower_id INTEGER, rate INTEGER, principal_amount INTEGER, state TEXT)")
	db.Exec(`CREATE TABLE loan_approvals (loan_id UUID, employee_id INTEGER, proof TEXT)`)
	db.Exec(`CREATE TABLE investments (
		id SERIAL PRIMARY KEY,
		loan_id UUID REFERENCES loans(id),
		amount INTEGER,
		investor_id INT
	);`)

	return db, nil
}
