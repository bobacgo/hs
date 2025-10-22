package hs

import (
	"os"
	"syscall"
	"time"
)

type Config struct {
	Addr            string
	ShutdownTimeout time.Duration
	sigs            []os.Signal
}

var defaultConfig = Config{
	ShutdownTimeout: 5 * time.Second,
	sigs:            []os.Signal{syscall.SIGINT, syscall.SIGKILL},
}

func WithShutdownTimeout(timeout time.Duration) func(*Config) {
	return func(c *Config) {
		c.ShutdownTimeout = timeout
	}
}

func WithSignals(sigs ...os.Signal) func(*Config) {
	return func(c *Config) {
		c.sigs = sigs
	}
}