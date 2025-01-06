package repository

import (
	"context"
	"nqrm/wbtechlvl0/order_services/internal/model"
)

type OrderDB interface {
	GetAllOrders(ctx context.Context) ([]model.Order, error)
	AddOrder(ctx context.Context, order model.Order) error
}

type CacheOrder interface {
	Get(orderUID string) (*model.Order, bool)
	Set(order *model.Order)
}
