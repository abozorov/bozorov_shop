package handlers

import (
	"net/http"

	"github.com/abozorov/bozorov_shop/internal/handlers/middleware"
	orderhandler "github.com/abozorov/bozorov_shop/internal/handlers/order"
	userhandler "github.com/abozorov/bozorov_shop/internal/handlers/user"
)

type Router struct {
	http.Handler
}

func NewRouter(u *userhandler.UserHandler, o *orderhandler.OrderHandler) *Router {


	mux := http.NewServeMux()

	// ---------- Public ----------
	mux.HandleFunc("POST /auth/register", u.Register)
	mux.HandleFunc("POST /auth/login", u.Login)

	
	return &Router{
		middleware.Logging(mux),
	}
}

