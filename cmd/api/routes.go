package main

import (
	"net/http"
	httpSwagger "github.com/swaggo/http-swagger"
	
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

		mux.Handle(
		"GET /swagger/",
		httpSwagger.WrapHandler,
	)

	mux.HandleFunc("GET /health", app.health)
	mux.HandleFunc("GET /menu", app.listDishes)
	mux.HandleFunc("GET /menu/{id}", app.showDish)
	mux.HandleFunc("POST /orders", app.createOrder)
	mux.HandleFunc("GET /orders", app.listOrders)
	mux.HandleFunc("GET /orders/{id}", app.showOrder)

	return mux
}
