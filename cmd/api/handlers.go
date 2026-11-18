package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	
    "github.com/binyam-ta/addis-eats-api/internal/delivery"
	"github.com/binyam-ta/addis-eats-api/internal/menu"
	"github.com/binyam-ta/addis-eats-api/internal/order"
)

type createOrderRequest struct {
	Area    string   `json:"area"`
	Phone   string   `json:"phone"`
	DishIDs []string `json:"dish_ids"`
}

type orderResponse struct {
	ID                string      `json:"id,omitempty"`
	CreatedAt         string      `json:"created_at,omitempty"`
	Items             []menu.Dish `json:"items"`
	Area              string      `json:"area"`
	Phone             string      `json:"phone,omitempty"`
	SubtotalSantim    int         `json:"subtotal_santim"`
	ServiceFeeSantim  int         `json:"service_fee_santim"`
	DeliveryFeeSantim int         `json:"delivery_fee_santim"`
	TotalSantim       int         `json:"total_santim"`
}


 // HealthCheck godoc
 // @Summary Check API health
 // @Description Check whether the Addis Eats API is running
 // @Tags Health
 // @Produce json
 // @Success 200 {object} envelope
 // @Router /health [get]
func (app *application) health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, envelope{
		"data": map[string]string{
			"status":      "available",
			"environment": "development",
		},
	})
}

// ListDishes godoc
// @Summary List all dishes
// @Description Get all available dishes from the menu
// @Tags Menu
// @Produce json
// @Success 200 {object} envelope
// @Failure 500 {object} envelope
// @Router /menu [get]
func (app *application) listDishes(w http.ResponseWriter, r *http.Request) {
	dishes, err := app.dishes.All()
	if err != nil {
		app.serverError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, envelope{"data": dishes})
}

// ShowDish godoc
// @Summary Get a dish
// @Description Get a single dish by its ID
// @Tags Menu
// @Produce json
// @Param id path string true "Dish ID"
// @Success 200 {object} envelope
// @Failure 404 {object} envelope
// @Failure 500 {object} envelope
// @Router /menu/{id} [get]
func (app *application) showDish(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	dish, err := app.dishes.ByID(id)
	switch {
	case errors.Is(err, menu.ErrNotFound):
		app.notFound(w, "dish not found")
	case err != nil:
		app.serverError(w, err)
	default:
		writeJSON(w, http.StatusOK, envelope{"data": dish})
	}
}

// CreateOrder godoc
// @Summary Create a new order
// @Description Create an order using selected menu dishes
// @Tags Orders
// @Accept json
// @Produce json
// @Param request body createOrderRequest true "Order details"
// @Success 201 {object} envelope
// @Failure 400 {object} envelope
// @Failure 500 {object} envelope
// @Router /orders [post]
func (app *application) createOrder(w http.ResponseWriter, r *http.Request) {
   
	// Limit body payload to 1MB
	r.Body = http.MaxBytesReader(w, r.Body, 1_048_576)

	var in createOrderRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&in); err != nil {
		app.badRequest(w, "malformed JSON body or unknown field")
		return
	}

	if len(in.DishIDs) == 0 {
		app.badRequest(w, "dish_ids must not be empty")
		return
	}

	// Lookup dishes from store to guarantee server-controlled pricing
	items := make([]menu.Dish, 0, len(in.DishIDs))
	for _, id := range in.DishIDs {
		d, err := app.dishes.ByID(id)
		if err != nil {
			app.badRequest(w, fmt.Sprintf("unknown dish %s", id))
			return
		}
		items = append(items, d)
	}

	o := order.New(in.Area, in.Phone, items)
	_, err := o.ComputeTotal()
	if err != nil {
		if errors.Is(err, delivery.ErrUnservedArea) {
			app.badRequest(w, fmt.Sprintf("we do not deliver to %s", in.Area))
			return
		}
		app.badRequest(w, err.Error())
		return
	}

	if err := app.orders.Create(&o); err != nil {
		app.serverError(w, err)
		return
	}

	// Set Location header prior to writing status code
	w.Header().Set("Location", "/orders/"+o.ID)
	writeJSON(w, http.StatusCreated, envelope{"data": o})
}

// ListOrders godoc
// @Summary List all orders
// @Description Retrieve all customer orders
// @Tags Orders
// @Produce json
// @Success 200 {object} envelope
// @Failure 500 {object} envelope
// @Router /orders [get]
func (app *application) listOrders(w http.ResponseWriter, r *http.Request) {
    orders, err := app.orders.All()
    if err != nil {
        app.serverError(w, err)
        return
    }

    writeJSON(w, http.StatusOK, envelope{
        "data": orders,
    })
}

// ShowOrder godoc
// @Summary Get an order
// @Description Retrieve a specific order using its ID
// @Tags Orders
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} envelope
// @Failure 404 {object} envelope
// @Failure 500 {object} envelope
// @Router /orders/{id} [get]
func (app *application) showOrder(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")

    ord, err := app.orders.ByID(id)

    switch {
    case errors.Is(err, order.ErrOrderNotFound):
        app.notFound(w, "order not found")
        return

    case err != nil:
        app.serverError(w, err)
        return

    default:
        writeJSON(w, http.StatusOK, envelope{
            "data": ord,
        })
    }
}