package app

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/Dzetner/lab-inventory-service/internal/db"
)

func registerChemicalRoutes(r *chi.Mux, a *App) {
	r.Get("/chemicals", a.listChemicalsHandler)
	r.Post("/chemicals", a.createChemicalHandler)
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
