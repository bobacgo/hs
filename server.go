package hs

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
)

type Engine struct {
	config  *Config
	srv     *http.Server
	handler http.Handler
}

func New(addr string, opts ...func(*Config)) *Engine {
	c := defaultConfig
	c.Addr = addr
	for _, opt := range opts {
		opt(&c)
	}
	return &Engine{
		config:  &c,
		handler: http.DefaultServeMux,
	}
}

func (e *Engine) SetHandler(handler http.Handler) {
	e.handler = handler
}

func (e *Engine) Run() error {
	e.srv = &http.Server{
		Addr:    e.config.Addr,
		Handler: e.handler,
	}

	go func() {
		slog.Info("Starting server", "address", e.config.Addr)
		if err := e.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Server failed to start", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, e.config.sigs...)
	<-quit

	return e.shutdown()
}

func (e *Engine) shutdown() error {
	slog.Info("Shutting down server gracefully")

	ctx, cancel := context.WithTimeout(context.Background(), e.config.ShutdownTimeout)
	defer cancel()

	if err := e.srv.Shutdown(ctx); err != nil {
		slog.Error("Server shutdown failed", "error", err)
		return err
	}
	slog.Info("Server stopped gracefully")
	return nil
}