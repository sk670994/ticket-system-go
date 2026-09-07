package main

import (
	"log"
	"net/http"
	"os"

	"ticket-system/handlers"
	"ticket-system/middleware"
	"ticket-system/store"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, _ = w.Write([]byte(`{"status":"ok"}`))
}

func main() {
	// Create the application's in-memory store.
	appStore := store.NewStore()

	// Create handlers using the shared store.
	authHandler := &handlers.AuthHandler{
		Store: appStore,
	}

	ticketHandler := &handlers.TicketHandler{
		Store: appStore,
	}

	// Create the HTTP router.
	mux := http.NewServeMux()

	// Public endpoints.
	mux.HandleFunc("GET /health", healthHandler)
	mux.HandleFunc("POST /auth/register", authHandler.Register)
	mux.HandleFunc("POST /auth/login", authHandler.Login)

	// Protected ticket endpoints.
	mux.Handle(
		"POST /tickets",
		middleware.Auth(http.HandlerFunc(ticketHandler.CreateTicket)),
	)

	mux.Handle(
		"GET /tickets",
		middleware.Auth(http.HandlerFunc(ticketHandler.ListTickets)),
	)

	mux.Handle(
		"GET /tickets/{id}",
		middleware.Auth(http.HandlerFunc(ticketHandler.GetTicket)),
	)

	mux.Handle(
		"PATCH /tickets/{id}/status",
		middleware.Auth(http.HandlerFunc(ticketHandler.UpdateTicketStatus)),
	)

	// Serve the frontend.
	mux.Handle("/", http.FileServer(http.Dir("./frontend")))

	// Use PORT from the environment for deployment.
	// Default to 8080 for local development.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on :%s", port)

	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}
