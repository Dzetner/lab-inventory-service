package app

import (
	"context"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Dzetner/lab-inventory-service/internal/db"
)

type App struct {
	Router *chi.Mux
	DB     *db.Queries
	Pool   *pgxpool.Pool
}

func New() (*App, error) {
	ctx := context.Background()

	dsn := os.Getenv("LAB_INVENTORY_DB_DSN")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5433/lab_inventory?sslmode=disable"
	}

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}

	a := &App{
		Router: chi.NewRouter(),
		DB:     db.New(pool),
		Pool:   pool,
	}

	a.routes()

	return a, nil
}

func (a *App) routes() {
	r := a.Router

	r.Get("/health", a.healthHandler)

	registerEmployeeRoutes(r, a)
	registerRoomRoutes(r, a)
	registerChemicalRoutes(r, a)
	registerContainerRoutes(r, a)
}

func (a *App) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}
