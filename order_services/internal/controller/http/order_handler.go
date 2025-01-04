package http

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func (hc *httpContoller) OrderHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderUID, ok := vars["order"]
	if !ok {
		http.Error(w, `{"order_details":""}`, http.StatusNotFound)
		return
	}
	order, err := hc.orderService.GetOrderByID(r.Context(), orderUID)

	if err != nil {
		http.Error(w, `{"order_details":"Order not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}
