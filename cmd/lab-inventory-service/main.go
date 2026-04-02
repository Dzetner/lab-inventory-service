package main

import (
	"log"
	"net/http"

	"github.com/Dzetner/lab-inventory-service/internal/app"
)

func main() {
	a, err := app.New()
	if err != nil {
		log.Fatalf("init app: %v", err)
	}
	defer a.Pool.Close()

	addr := ":8080"
	log.Println("listening on", addr)

	if err := http.ListenAndServe(addr, a.Router); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
