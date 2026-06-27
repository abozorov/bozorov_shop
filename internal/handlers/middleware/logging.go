package middleware

import (
	"log"
	"net/http"
	"time"
)

func (m *Middleware)Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		startTime := time.Now()
		log.Printf("[INFO]	START http.method: {%s} url: {%s}\n",
			r.Method,
			r.URL.Path,
		)

		next.ServeHTTP(w, r)

		endTime := time.Now()
		log.Printf("[INFO]	END request duration: {%d} ms\n\n",
			int(endTime.UnixMilli())-int(startTime.UnixMilli()))
	})
}
