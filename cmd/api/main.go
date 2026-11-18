package main

import (
	"log"
	"net/http"
	"os"
	"time"

    _"github.com/binyam-ta/addis-eats-api/docs"
	"github.com/binyam-ta/addis-eats-api/internal/menu"
	"github.com/binyam-ta/addis-eats-api/internal/order"
)

type application struct {
	logger *log.Logger
	dishes menu.DishReader
	orders order.Store
}

// @title Addis Eats API
// @version 1.0
// @description REST API documentation for the Addis Eats food ordering platform.
// @host localhost:8080
// @BasePath /
// @schemes http
func main() {
	logger := log.New(os.Stdout, "", log.LstdFlags)

	app := &application{
		logger: logger,
		dishes: menu.NewMemoryStore(),
		orders: order.NewMemoryStore(),
	}

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      app.routes(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  time.Minute,
	}

	logger.Printf("starting server on %s", srv.Addr)
	logger.Fatal(srv.ListenAndServe())



}