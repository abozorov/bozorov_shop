package middleware

import (
	"context"
	"net/http"
	"strings"

	mycontext "github.com/abozorov/bozorov_shop/internal/my_context"
	"github.com/abozorov/bozorov_shop/pkg/jwt"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")

		if !strings.HasPrefix(token, "Bearer ") {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		token = strings.TrimPrefix(token, "Bearer ")
		claims, err := jwt.ParseToken(token)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), mycontext.UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, mycontext.EmailKey, claims.Email)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
