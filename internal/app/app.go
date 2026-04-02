package app

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
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

	r.Get("/employees", a.listEmployeesHandler)
	r.Post("/employees", a.createEmployeeHandler)

	r.Get("/rooms", a.listRoomsHandler)
	r.Post("/rooms", a.createRoomHandler)

	r.Get("/chemicals", a.listChemicalsHandler)
	r.Post("/chemicals", a.createChemicalHandler)
}

func (a *App) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

func (a *App) listEmployeesHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	emps, err := a.DB.ListEmployees(ctx)
	if err != nil {
		log.Println("list employees:", err)
		http.Error(w, "failed to list employees", http.StatusInternalServerError)
		return
	}

	if emps == nil {
		emps = []db.Employee{}
	}

	writeJSON(w, http.StatusOK, emps)
}

func (a *App) createEmployeeHandler(w http.ResponseWriter, r *http.Request) {
	type input struct {
		FullName string `json:"full_name"`
		Role     string `json:"role"`
	}

	var in input
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if in.FullName == "" || in.Role == "" {
		http.Error(w, "full_name and role are required", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	emp, err := a.DB.CreateEmployee(ctx, db.CreateEmployeeParams{
		FullName: in.FullName,
		Role:     in.Role,
	})
	if err != nil {
		log.Println("create employee:", err)
		http.Error(w, "failed to create employee", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, emp)
}

func (a *App) listRoomsHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	rooms, err := a.DB.ListRooms(ctx)
	if err != nil {
		log.Println("list rooms:", err)
		http.Error(w, "failed to list rooms", http.StatusInternalServerError)
		return
	}

	if rooms == nil {
		rooms = []db.Room{}
	}

	writeJSON(w, http.StatusOK, rooms)
}

func (a *App) createRoomHandler(w http.ResponseWriter, r *http.Request) {
	type input struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	var in input
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if in.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	room, err := a.DB.CreateRoom(ctx, db.CreateRoomParams{
		Name: in.Name,
		Description: pgtype.Text{
			String: in.Description,
			Valid:  in.Description != "",
		},
	})
	if err != nil {
		log.Println("create room:", err)
		http.Error(w, "failed to create room", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, room)
}

func (a *App) listChemicalsHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	chems, err := a.DB.ListChemicals(ctx)
	if err != nil {
		log.Println("list chemicals:", err)
		http.Error(w, "failed to list chemicals", http.StatusInternalServerError)
		return
	}

	if chems == nil {
		chems = []db.Chemical{}
	}

	writeJSON(w, http.StatusOK, chems)
}

func (a *App) createChemicalHandler(w http.ResponseWriter, r *http.Request) {
	type input struct {
		Name      string `json:"name"`
		CASNumber string `json:"cas_number"`
		Formula   string `json:"formula"`
		SDSURL    string `json:"sds_url"`
	}

	var in input
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if in.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	chem, err := a.DB.CreateChemical(ctx, db.CreateChemicalParams{
		Name: in.Name,
		CasNumber: pgtype.Text{
			String: in.CASNumber,
			Valid:  in.CASNumber != "",
		},
		Formula: pgtype.Text{
			String: in.Formula,
			Valid:  in.Formula != "",
		},
		SdsUrl: pgtype.Text{
			String: in.SDSURL,
			Valid:  in.SDSURL != "",
		},
	})
	if err != nil {
		log.Println("create chemical:", err)
		http.Error(w, "failed to create chemical", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, chem)
}
