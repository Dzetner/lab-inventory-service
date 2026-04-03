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

func registerRoomRoutes(r *chi.Mux, a *App) {
	r.Get("/rooms", a.listRoomsHandler)
	r.Post("/rooms", a.createRoomHandler)
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
