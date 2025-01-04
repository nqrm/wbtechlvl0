package services

import (
	"context"
	"errors"
	"fmt"
	"nqrm/wbtechlvl0/order_services/internal/model"
	"nqrm/wbtechlvl0/order_services/internal/repository"
)

type OrderService struct {
	db    repository.OrderDB
	cache repository.CacheOrder
}

func NewOrderService(db repository.OrderDB, cache repository.CacheOrder) *OrderService {
	return &OrderService{db, cache}
}

func (os *OrderService) GetOrderByID(ctx context.Context, orderUID string) (model.Order, error) {
	order, ok := os.cache.Get(orderUID)
	if !ok {
		orderNotFoundError := fmt.Sprintf("Order with UID %v not found", orderUID)
		return model.Order{}, errors.New(orderNotFoundError)
	}

	return *order, nil
}
