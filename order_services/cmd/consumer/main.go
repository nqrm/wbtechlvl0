package main

import "nqrm/wbtechlvl0/order_services/internal/app"

func main() {
	/*
		query := `SELECT * FROM orders WHERE order_details->>'order_uid' = $1`
		var orderData string
		var orderID string = "b563feb7b2b84b6test"
		err = conn.QueryRow(context.Background(), query, orderID).Scan(&orderData)
		if err != nil {
			log.Fatalf("Query failed: %v\n", err)
		}*/
	app.Run()
}
