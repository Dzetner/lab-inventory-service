package app

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/Dzetner/lab-inventory-service/internal/db"
)

func registerEmployeeRoutes(r *chi.Mux, a *App) {
	r.Get("/employees", a.listEmployeesHandler)
	r.Post("/employees", a.createEmployeeHandler)
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
