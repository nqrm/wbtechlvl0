package pgrepo

import (
	"context"
	"encoding/json"
	"log"
	"nqrm/wbtechlvl0/order_services/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type postgres struct {
	db *pgxpool.Pool
}

func NewPG(db *pgxpool.Pool) *postgres {
	return &postgres{db}
}

// вынести десериализацию, возвращать []byte
// в запросе выбирать конкретные поля, а не *
func (pg *postgres) GetOrderByID(ctx context.Context, orderUID string) (model.Order, error) {
	query := `SELECT * FROM orders WHERE order_details->>'order_uid' = $1`

	var orderData string
	err := pg.db.QueryRow(ctx, query, orderUID).Scan(&orderData)
	if err != nil {
		log.Printf("Query failed: %v\n", err)
		return model.Order{}, err
	}

	var order model.Order
	err = json.Unmarshal([]byte(orderData), &order)
	if err != nil {
		log.Printf("Failed to unmarshal JSON: %v\n", err)
		return model.Order{}, err
	}

	return order, err
}

// вынести десериализацию, возвращать [][]byte
func (pg *postgres) GetAllOrders(ctx context.Context) ([]model.Order, error) {
	query := `SELECT order_details FROM orders`

	rows, err := pg.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	var orders []model.Order
	for rows.Next() {
		var order model.Order
		var orderData string
		if err := rows.Scan(&orderData); err != nil {
			return nil, err
		}
		err = json.Unmarshal([]byte(orderData), &order)
		if err != nil {
			log.Fatalf("Failed to unmarshal JSON: %v\n", err)
			return nil, err
		}
		orders = append(orders, order)

	}

	return orders, nil
}

// вынести сериализацию, передавать []byte
func (pg *postgres) AddOrder(ctx context.Context, order model.Order) error {

	orderJSON, err := json.Marshal(order)
	if err != nil {
		log.Fatalf("Failed to marshal order to JSON: %v\n", err)
	}

	query := `INSERT INTO orders (order_details) VALUES ($1)`

	_, err = pg.db.Exec(context.Background(), query, orderJSON)
	if err != nil {
		log.Fatalf("Failed to insert order into database: %v\n", err)
	}

	return err
}
