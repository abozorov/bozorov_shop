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

	// ---------- Public ----------
	mux.HandleFunc("POST /auth/register", u.Register)
	mux.HandleFunc("POST /auth/login", u.Login)
	mux.HandleFunc("POST /user/verify", u.Verify)

	// ---------- Authorized user ----------
	mux.Handle("GET /users/me", //✅
		middleware.Auth(http.HandlerFunc(u.GetMe)),
	)
	mux.Handle("PATCH /users/me", //✅
		middleware.Auth(http.HandlerFunc(u.UpdateMe)),
	)
	mux.Handle("DELETE /users/me", //✅
		middleware.Auth(http.HandlerFunc(u.DeleteMe)),
	)

	mux.Handle("POST /orders", //✅
		middleware.Auth(http.HandlerFunc(o.CreateOrder)),
	)
	mux.Handle("GET /orders", //✅
		middleware.Auth(http.HandlerFunc(o.GetMyOrders)),
	)
	mux.Handle("GET /orders/{id}", //✅
		middleware.Auth(http.HandlerFunc(o.GetOrderByID)),
	)
	mux.Handle("PATCH /orders/{id}", //✅
		middleware.Auth(http.HandlerFunc(o.UpdateOrder)),
	)
	mux.Handle("DELETE /orders/{id}", //✅
		middleware.Auth(http.HandlerFunc(o.CancleOrder)),
	)

	// ---------- Admin ----------
	mux.Handle("GET /admin/users", //✅ (Admin)
		middleware.AuthAdmin(http.HandlerFunc(u.GetAllUsers)),
	)
	mux.Handle("GET /admin/orders", //✅ (Admin)
		middleware.AuthAdmin(http.HandlerFunc(o.GetAllOrders)),
	)
	mux.Handle("PATCH /admin/users/{id}/role", //✅ (Admin)
		middleware.AuthAdmin(http.HandlerFunc(u.UpdateUserRole)),
	)

	// Logging применяется ко всем маршрутам, включая public.
	return &Router{
		middleware.Logging(mux),
	}
}

/*

	| Метод  | Endpoint                 | Авторизация  | Описание                    |
	| ------ | ------------------------ | ------------ | --------------------------- |
	| POST   | `/auth/register`         | ❌           | Регистрация                 |
	| POST   | `/auth/login`            | ❌           | Авторизация                 |
	| GET    | `/users/me`              | ✅           | Получить свой профиль       |
	| PUT    | `/users/me`              | ✅           | Обновить профиль            |
	| DELETE | `/users/me`              | ✅           | Удалить аккаунт             |
	| POST   | `/orders`                | ✅           | Создать заказ               |
	| GET    | `/orders`                | ✅           | Получить свои заказы        |
	| GET    | `/orders/{id}`           | ✅           | Получить заказ по ID        |
	| PUT    | `/orders/{id}`           | ✅           | Обновить заказ              |
	| DELETE | `/orders/{id}`           | ✅           | Удалить заказ               |
	| GET    | `/admin/users`           | ✅ (Admin)   | Получить всех пользователей |
	| GET    | `/admin/orders`          | ✅ (Admin)   | Получить все заказы         |
	| PATCH  | `/admin/users/{id}/role` | ✅ (Admin)   | Изменить роль               |

*/
