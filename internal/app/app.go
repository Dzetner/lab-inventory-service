package app

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
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

	r.Get("/containers", a.listContainersHandler)
	r.Post("/containers", a.createContainerHandler)
	r.Post("/containers/{id}/checkout", a.checkoutContainerHandler)
	r.Post("/containers/{id}/return", a.returnContainerHandler)
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

func (a *App) listContainersHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	containers, err := a.DB.ListContainers(ctx)
	if err != nil {
		log.Println("list containers:", err)
		http.Error(w, "failed to list containers", http.StatusInternalServerError)
		return
	}

	if containers == nil {
		containers = []db.Container{}
	}

	writeJSON(w, http.StatusOK, containers)
}

func (a *App) createContainerHandler(w http.ResponseWriter, r *http.Request) {
	type input struct {
		ChemicalID   int64   `json:"chemical_id"`
		RoomID       int64   `json:"room_id"`
		LabelCode    string  `json:"label_code"`
		Quantity     float64 `json:"quantity"`
		Unit         string  `json:"unit"`
		Status       string  `json:"status"`
		CheckedOutBy *int64  `json:"checked_out_by"`
	}

	var in input
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if in.ChemicalID == 0 || in.RoomID == 0 || in.Unit == "" {
		http.Error(w, "chemical_id, room_id and unit are required", http.StatusBadRequest)
		return
	}
	if in.Status == "" {
		in.Status = "available"
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	labelCode := pgtype.Text{
		String: in.LabelCode,
		Valid:  in.LabelCode != "",
	}

	var checkedOutBy pgtype.Int8
	if in.CheckedOutBy != nil {
		checkedOutBy = pgtype.Int8{
			Int64: *in.CheckedOutBy,
			Valid: true,
		}
	} else {
		checkedOutBy = pgtype.Int8{Valid: false}
	}

	container, err := a.DB.CreateContainer(ctx, db.CreateContainerParams{
		ChemicalID:   in.ChemicalID,
		RoomID:       in.RoomID,
		LabelCode:    labelCode,
		Quantity:     in.Quantity,
		Unit:         in.Unit,
		Status:       in.Status,
		CheckedOutBy: checkedOutBy,
	})
	if err != nil {
		log.Println("create container:", err)
		http.Error(w, "failed to create container", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, container)
}

func (a *App) checkoutContainerHandler(w http.ResponseWriter, r *http.Request) {
	type input struct {
		EmployeeID int64 `json:"employee_id"`
	}

	idStr := chi.URLParam(r, "id")
	containerID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || containerID <= 0 {
		http.Error(w, "invalid container id", http.StatusBadRequest)
		return
	}

	var in input
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if in.EmployeeID == 0 {
		http.Error(w, "employee_id is required", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	checkedOutBy := pgtype.Int8{
		Int64: in.EmployeeID,
		Valid: true,
	}

	container, err := a.DB.CheckoutContainer(ctx, db.CheckoutContainerParams{
		ID:           containerID,
		CheckedOutBy: checkedOutBy,
	})
	if err != nil {
		log.Println("checkout container:", err)
		http.Error(w, "failed to checkout container", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, container)
}

func (a *App) returnContainerHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	containerID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || containerID <= 0 {
		http.Error(w, "invalid container id", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	container, err := a.DB.ReturnContainer(ctx, containerID)
	if err != nil {
		log.Println("return container:", err)
		http.Error(w, "failed to return container", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, container)
}
