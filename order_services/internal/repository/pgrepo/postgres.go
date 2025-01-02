package pgrepo

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"nqrm/wbtechlvl0/order_services/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type postgres struct {
	db *pgxpool.Pool
}

func NewPG(connStr string) (*postgres, error) {
	db, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return nil, err
	}

	return &postgres{db: db}, nil
}

func (pg *postgres) Ping(ctx context.Context) error {
	return pg.db.Ping(ctx)
}

func (pg *postgres) Close() {
	pg.db.Close()
}

func (pg *postgres) GetOrderByID(ctx context.Context, orderUID string) (model.Order, error) {
	query := `SELECT * FROM orders WHERE order_details->>'order_uid' = $1`

	var orderData string
	err := pg.db.QueryRow(ctx, query, orderUID).Scan(&orderData)
	if err != nil {
		log.Fatalf("Query failed: %v\n", err)
	}

	var order model.Order
	err = json.Unmarshal([]byte(orderData), &order)
	if err != nil {
		log.Fatalf("Failed to unmarshal JSON: %v\n", err)
	}

	fmt.Printf("Parsed Order: %+v\n", order)

	return order, nil
}
