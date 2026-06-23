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
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type responseOrder struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
}

func newResponseOrder(u models.Order) *responseOrder {
	return &responseOrder{
		CreatedAt: u.CreatedAt.Format(time.RFC822Z),
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Phone:     u.Phone,
		Role:      u.Role,
	}
}

func (h *OrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// get order
	usr := requestOrder{}
	err := json.NewDecoder(r.Body).Decode(&usr)
	if err != nil {
		errs.ErrsToHttp(w, errs.ErrBadRequestBody)
		return
	}

	// creating & transform models.Order -> order
	err = h.service.Create(r.Context(), models.Order{
		Name:         usr.Name,
		Email:        usr.Email,
		Phone:        usr.Phone,
		PasswordHash: usr.Password,
		Role:         usr.Role,
	})
	if err != nil {
		h.logger.Error("order_handler.Create: ", zap.String("error", err.Error()))
		errs.ErrsToHttp(w, err)
		return
	}
	w.Write([]byte("Order Created"))
}

func (h *OrderHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// load all
	ctx, cancle := context.WithTimeout(r.Context(), time.Second*2)
	defer cancle()
	orders, err := h.service.GetAll(ctx)
	if err != nil {
		h.logger.Error("order_handler.GetAll: ", zap.String("error", err.Error()))
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
		h.logger.Error("order_handler.GetAll: ", zap.String("error", err.Error()))
		errs.ErrsToHttp(w, err)
		return
	}
}

func (h *OrderHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// check path
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		h.logger.Error("order_handler.GetByID: ", zap.String("error", err.Error()))
		errs.ErrsToHttp(w, errs.ErrBadRequest)
		return
	}

	// get by id
	usr, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		h.logger.Error("order_handler.GetByID: ", zap.String("error", err.Error()))
		errs.ErrsToHttp(w, err)
		return
	}

	// transform models.Order -> order
	resp := *newResponseOrder(*usr)

	// write response
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		h.logger.Error("order_handler.GetAll: ", zap.String("error", err.Error()))
		errs.ErrsToHttp(w, err)
	}
}

func (h *OrderHandler) Update(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// get order
	usr := requestOrder{}
	err := json.NewDecoder(r.Body).Decode(&usr)
	if err != nil {
		errs.ErrsToHttp(w, errs.ErrBadRequestBody)
		return
	}

	// creating & transform models.Order -> order
	err = h.service.Update(r.Context(), models.Order{
		ID:    usr.ID,
		Name:  usr.Name,
		Email: usr.Email,
		Phone: usr.Phone,
		Role:  usr.Role,
	})
	if err != nil {
		h.logger.Error("order_handler.Update: ", zap.String("error", err.Error()))
		errs.ErrsToHttp(w, err)
		return
	}
	w.Write([]byte("Order updated"))
}

func (h *OrderHandler) Delete(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// check path
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		h.logger.Error("order_handler.Delete: ", zap.String("error", err.Error()))
		errs.ErrsToHttp(w, errs.ErrBadRequest)
		return
	}

	// get by id
	err = h.service.DeleteOrder(r.Context(), id)
	if err != nil {
		h.logger.Error("order_handler.Delete: ", zap.String("error", err.Error()))
		errs.ErrsToHttp(w, err)
		return
	}

	// write response
	w.Write([]byte(fmt.Sprintf("order %d deleted", id)))
}
