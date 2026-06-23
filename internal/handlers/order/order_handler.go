package orderhandler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/abozorov/bozorov_shop/internal/models"
	orderservice "github.com/abozorov/bozorov_shop/internal/service/order"
	"github.com/abozorov/bozorov_shop/logger"
	"github.com/abozorov/bozorov_shop/pkg/errs"
	"go.uber.org/zap"
)

type OrderHandler struct {
	service *orderservice.OrderService
	logger  *logger.Logger
}

func NewOrderHandler(service *orderservice.OrderService, logger *logger.Logger) *OrderHandler {
	return &OrderHandler{
		service: service,
		logger:  logger,
	}
}

type requestOrder struct {
	ID      int     `json:"id"`
	UserID  int     `json:"user_id"`
	Product string  `json:"product"`
	Price   float64 `json:"price"`
	Status  string  `json:"status"`
}

type responseOrder struct {
	ID        int     `json:"id"`
	UserID    int     `json:"user_id"`
	Product   string  `json:"product"`
	Price     float64 `json:"price"`
	Status    string  `json:"status"`
	CreatedAt string  `json:"created_at"`
}

func newResponseOrder(o models.Order) *responseOrder {
	return &responseOrder{
		ID:        o.ID,
		UserID:    o.UserID,
		Product:   o.Product,
		Price:     o.Price,
		Status:    o.Status,
		CreatedAt: o.CreatedAt.Format(time.RFC822Z),
	}
}

func (o *OrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// get order
	order := requestOrder{}
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		errs.ErrsToHttp(w, errs.ErrBadRequestBody)
		return
	}

	// creating & transform models.Order -> order
	err = o.service.Create(r.Context(), models.Order{
		UserID:  order.UserID,
		Product: order.Product,
		Price:   order.Price,
	})
	if err != nil {
		o.logger.Error("order_handler.Create: ", zap.String("error", err.Error()))
		errs.ErrsToHttp(w, err)
		return
	}
	w.Write([]byte("Order Created"))
}

func (o *OrderHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// load all
	ctx, cancle := context.WithTimeout(r.Context(), time.Second*2)
	defer cancle()
	orders, err := o.service.GetAll(ctx)
	if err != nil {
		o.logger.Error("order_handler.GetAll: ", zap.String("error", err.Error()))
		errs.ErrsToHttp(w, err)
		return
	}

	// transform models.Order -> order
	resp := make([]responseOrder, 0, len(orders))
	for _, v := range orders {
		resp = append(resp, *newResponseOrder(v))
	}

	// write request
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		o.logger.Error("order_handler.GetAll: ", zap.String("error", err.Error()))
		errs.ErrsToHttp(w, err)
		return
	}
}

func (o *OrderHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// check path
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		o.logger.Error("order_handler.GetByID: ", zap.String("error", err.Error()))
		errs.ErrsToHttp(w, errs.ErrBadRequest)
		return
	}

	// get by id
	order, err := o.service.GetByID(r.Context(), id)
	if err != nil {
		o.logger.Error("order_handler.GetByID: ", zap.String("error", err.Error()))
		errs.ErrsToHttp(w, err)
		return
	}

	// transform models.Order -> order
	resp := *newResponseOrder(*order)

	// write response
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		o.logger.Error("order_handler.GetAll: ", zap.String("error", err.Error()))
		errs.ErrsToHttp(w, err)
	}
}

func (o *OrderHandler) Update(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// get order
	order := requestOrder{}
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		errs.ErrsToHttp(w, errs.ErrBadRequestBody)
		return
	}

	// creating & transform models.Order -> order
	err = o.service.Update(r.Context(), models.Order{
		ID:      order.ID,
		UserID:  order.UserID,
		Product: order.Product,
		Price:   order.Price,
		Status:  order.Status,
	})
	if err != nil {
		o.logger.Error("order_handler.Update: ", zap.String("error", err.Error()))
		errs.ErrsToHttp(w, err)
		return
	}
	w.Write([]byte("Order updated"))
}

func (o *OrderHandler) CancleOrder(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// check path
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		o.logger.Error("order_handler.CancleOrder: ", zap.String("error", err.Error()))
		errs.ErrsToHttp(w, errs.ErrBadRequest)
		return
	}

	// get by id
	err = o.service.CancleOrder(r.Context(), id)
	if err != nil {
		o.logger.Error("order_handler.CancleOrder: ", zap.String("error", err.Error()))
		errs.ErrsToHttp(w, err)
		return
	}

	// write response
	w.Write([]byte(fmt.Sprintf("order %d cancled", id)))
}
