package handlers

import (
	"net/http"

	"github.com/abozorov/bozorov_shop/internal/handlers/middleware"
	userhandler "github.com/abozorov/bozorov_shop/internal/handlers/user"
)

type Router struct {
	*http.ServeMux
}

func NeweRouter(u userhandler.UserHandler) *Router {
	mux := http.NewServeMux()

	// USER
	mux.Handle("POST /users", middleware.Logging(http.HandlerFunc(u.Create)))
	mux.Handle("GET /users", middleware.Logging(http.HandlerFunc(u.GetAll)))
	mux.Handle("GET /user/{id}", middleware.Logging(http.HandlerFunc(u.GetByID)))
	mux.Handle("PUT /user", middleware.Logging(http.HandlerFunc(u.Update)))
	mux.Handle("DELETE /user/{id}", middleware.Logging(http.HandlerFunc(u.Delete)))

	return &Router{
		mux,
	}
}
