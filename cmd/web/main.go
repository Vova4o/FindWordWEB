package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/Vova4o/findaWordweb/config"
)

type application struct {
	cfg         config.Config
	templateMap map[string]*template.Template
}

func main() {
	flag.Parse()

	app := application{
		cfg:         config.NewConfig(),
		templateMap: make(map[string]*template.Template),
	}

	// Initialize the router
	router := app.routes()

	// Combine host and port for the full address
	address := app.cfg.Host + app.cfg.Port

	// Create a new http.Server instance
	server := &http.Server{
		Addr:              address,
		Handler:           router,
		IdleTimeout:       30 * time.Second,
		ReadTimeout:       30 * time.Second, // Set the ReadTimeout to 5 seconds
		ReadHeaderTimeout: 30 * time.Second, // Set the ReadHeaderTimeout to 2 seconds
		WriteTimeout:      30 * time.Second, // Set the WriteTimeout to 10 seconds
	}

	// Start the server
	fmt.Println("Starting server on", address)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
