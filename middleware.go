package hs

import (
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				slog.ErrorContext(r.Context(), "Recovered from panic", slog.Any("error", err))
				debug.PrintStack()
				JSON(w, map[string]any{
					"code": 1,
					"msg":  "Internal Server Error",
				})
			}
		}()

		now := time.Now()
		next.ServeHTTP(w, r)
		slog.InfoContext(r.Context(), "Request received",
			slog.String("time", time.Since(now).String()),
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.String("remote_addr", r.RemoteAddr),
		)
	})
}