package handlers

import (
	"net/http"

	"github.com/abozorov/bozorov_shop/internal/handlers/middleware"
	orderhandler "github.com/abozorov/bozorov_shop/internal/handlers/order"
	userhandler "github.com/abozorov/bozorov_shop/internal/handlers/user"
)

type Router struct {
	*http.ServeMux
}

func NewRouter(u *userhandler.UserHandler, o *orderhandler.OrderHandler) *Router {
	mux := http.NewServeMux()

	// USER
	mux.Handle("POST /users", middleware.Logging(http.HandlerFunc(u.Create)))
	mux.Handle("GET /users", middleware.Logging(http.HandlerFunc(u.GetAll)))
	mux.Handle("GET /user/{id}", middleware.Logging(http.HandlerFunc(u.GetByID)))
	mux.Handle("PUT /user", middleware.Logging(http.HandlerFunc(u.Update)))
	mux.Handle("DELETE /user/{id}", middleware.Logging(http.HandlerFunc(u.Delete)))

	// ORDERS
	mux.Handle("POST /orders", middleware.Logging(http.HandlerFunc(o.Create)))
	mux.Handle("GET /orders", middleware.Logging(http.HandlerFunc(o.GetAll)))
	mux.Handle("GET /order/{id}", middleware.Logging(http.HandlerFunc(o.GetByID)))
	mux.Handle("PUT /order", middleware.Logging(http.HandlerFunc(o.Update)))
	mux.Handle("DELETE /order/{id}", middleware.Logging(http.HandlerFunc(o.CancleOrder)))

	return &Router{
		mux,
	}
}
