package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/abozorov/bozorov_shop/internal/models"
	mycontext "github.com/abozorov/bozorov_shop/internal/my_context"
	"github.com/abozorov/bozorov_shop/pkg/jwt"
)

type IMiddleware interface {
	GetByID(ctx context.Context, id int) (*models.User, error)
}

type Middleware struct {
	repo IMiddleware
	jwt  *jwt.JWTSecret
}

func NewMiddlware(repo IMiddleware, jwt *jwt.JWTSecret) *Middleware {
	return &Middleware{
		repo: repo,
		jwt:  jwt,
	}
}

func (m *Middleware) AuthAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")

		if !strings.HasPrefix(token, "Bearer ") {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		token = strings.TrimPrefix(token, "Bearer ")
		claims, err := m.jwt.ParseToken(token)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		user, err := m.repo.GetByID(r.Context(), claims.UserID)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if user.Role != "admin" || !user.DeletedAt.IsZero() {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), mycontext.UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, mycontext.EmailKey, claims.Email)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *Middleware) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")

		if !strings.HasPrefix(token, "Bearer ") {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		token = strings.TrimPrefix(token, "Bearer ")
		claims, err := m.jwt.ParseToken(token)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		user, err := m.repo.GetByID(r.Context(), claims.UserID)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		if !user.DeletedAt.IsZero() {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), mycontext.UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, mycontext.EmailKey, claims.Email)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
