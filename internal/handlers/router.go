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

func NewRouter(u *userhandler.UserHandler, o *orderhandler.OrderHandler, middleware *middleware.Middleware) *Router {
	mux := http.NewServeMux()

	// Public
	mux.HandleFunc("POST /auth/register", u.Register)
	mux.HandleFunc("POST /auth/login", u.Login)
	mux.HandleFunc("POST /user/verify", u.Verify)
	mux.HandleFunc("POST /auth/refresh", u.RefreshToken)

	// Authorized user: USERS
	mux.Handle("GET /users/me",
		middleware.Auth(http.HandlerFunc(u.GetMe)),
	)
	mux.Handle("GET /users/profile",
		middleware.Auth(http.HandlerFunc(u.GetProfile)),
	)
	mux.Handle("GET /users/login-history",
		middleware.Auth(http.HandlerFunc(u.LoginHistory)),
	)
	mux.Handle("PATCH /users/me",
		middleware.Auth(http.HandlerFunc(u.UpdateMe)),
	)
	mux.Handle("DELETE /users/me",
		middleware.Auth(http.HandlerFunc(u.DeleteMe)),
	)
	mux.Handle("PATCH /users/change-password",
		middleware.Auth(http.HandlerFunc(u.ChangePassword)),
	)
	mux.Handle("POST /auth/logout",
		middleware.Auth(http.HandlerFunc(u.Logout)),
	)

	// ORDERS
	mux.Handle("POST /orders",
		middleware.Auth(http.HandlerFunc(o.CreateOrder)),
	)
	mux.Handle("GET /orders",
		middleware.Auth(http.HandlerFunc(o.GetMyOrders)),
	)
	mux.Handle("GET /orders/{id}",
		middleware.Auth(http.HandlerFunc(o.GetOrderByID)),
	)
	mux.Handle("PATCH /orders/{id}",
		middleware.Auth(http.HandlerFunc(o.UpdateOrder)),
	)
	mux.Handle("PATCH /orders/{id}/cancle",
		middleware.Auth(http.HandlerFunc(o.CancleOrder)),
	)

	// Authorized user: Admin
	mux.Handle("GET /admin/users",
		middleware.AuthAdmin(http.HandlerFunc(u.GetAllUsers)),
	)
	mux.Handle("GET /admin/orders",
		middleware.AuthAdmin(http.HandlerFunc(o.GetAllOrders)),
	)
	mux.Handle("PATCH /admin/users/{id}/role",
		middleware.AuthAdmin(http.HandlerFunc(u.UpdateUserRole)),
	)

	// Logging
	return &Router{
		middleware.Logging(mux),
	}
}
