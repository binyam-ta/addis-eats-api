
package main

import (
	"fmt"
	"strings"
)

type Dish struct {
	Name  string
	Price int 
}

func main() {
	// Declare dishes
	dishes := []Dish{
		{Name: "Doro Wat", Price: 18500},
		{Name: "Shiro", Price: 12000},
		{Name: "Tibs", Price: 24000},
		{Name: "Injera basket", Price: 4500},
	}

	subtotal := 0

	fmt.Println("========== ADDIS EATS RECEIPT ==========")

	// Print each dish and calculate subtotal
	for _, dish := range dishes {
		fmt.Printf("%-20s ETB %.2f\n", dish.Name, float64(dish.Price)/100)
		subtotal += dish.Price
	}

	fmt.Println("----------------------------------------")
	fmt.Printf("%-20s ETB %.2f\n", "Subtotal:", float64(subtotal)/100)

	
	serviceFee := int(float64(subtotal) * 0.05)

	fmt.Printf("%-20s ETB %.2f\n", "Service fee:", float64(serviceFee)/100)

	deliveryArea := "Bole"

	deliveryFee := 0

	// Free delivery for orders reaching 50000 santim
	if subtotal >= 50000 {
		deliveryFee = 0
	} else {
		switch deliveryArea {
		case "Bole":
			deliveryFee = 3000
		case "Kazanchis":
			deliveryFee = 3500
		case "Piassa":
			deliveryFee = 4000
		case "Megenagna":
			deliveryFee = 4500
		default:
			deliveryFee = 5000
		}
	}

	fmt.Printf("%-20s ETB %.2f\n", "Delivery fee:", float64(deliveryFee)/100)

	grandTotal := subtotal + serviceFee + deliveryFee

	fmt.Println("----------------------------------------")
	fmt.Printf("%-20s ETB %.2f\n", "Grand total:", float64(grandTotal)/100)

	customerPhone := "0911223344"

	if strings.HasPrefix(customerPhone, "09") {
		fmt.Println("Customer phone:", customerPhone)
	}

	fmt.Println("========================================")
}
