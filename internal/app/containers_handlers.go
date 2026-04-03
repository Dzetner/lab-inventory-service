package app

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/Dzetner/lab-inventory-service/internal/db"
)

func registerContainerRoutes(r *chi.Mux, a *App) {
	r.Get("/containers", a.listContainersHandler)
	r.Post("/containers", a.createContainerHandler)
	r.Post("/containers/{id}/checkout", a.checkoutContainerHandler)
	r.Post("/containers/{id}/return", a.returnContainerHandler)
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
