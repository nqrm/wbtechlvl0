package http

import (
	"context"
	"net/http"
	"nqrm/wbtechlvl0/order_services/internal/model"

	"github.com/gorilla/mux"
)

type OrderService interface {
	GetOrderByID(ctx context.Context, orderUID string) (model.Order, error)
}

type httpContoller struct {
	orderService OrderService
}

func NewOrderRouter(orderService OrderService) *mux.Router {
	httpContoller := &httpContoller{orderService}
	r := mux.NewRouter()
	r.HandleFunc("/orders/{order}", httpContoller.OrderHandler).Methods(http.MethodGet)
	return r
}
