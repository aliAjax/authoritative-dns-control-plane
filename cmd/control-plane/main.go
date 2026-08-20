package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"example.com/authoritativedns/internal/api"
	"example.com/authoritativedns/internal/repository"
	"example.com/authoritativedns/internal/resolver_runtime"
	"example.com/authoritativedns/internal/service"
)

func main() {
	store := repository.NewMemoryStore()
	app := service.NewApplication(store)
	mux := api.NewRouter(app)
	httpSrv := &http.Server{Addr: env("HTTP_ADDR", ":8080"), Handler: mux, ReadHeaderTimeout: 5 * time.Second}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	go func() {
		log.Printf("http listening on %s", httpSrv.Addr)
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("http: %v", err)
		}
	}()
	resolver := resolver_runtime.NewServer(app)
	go func() {
		if err := resolver.ListenAndServeUDP(ctx, env("DNS_ADDR", ":8053")); err != nil {
			log.Printf("dns udp: %v", err)
		}
	}()
	<-ctx.Done()
	shutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = httpSrv.Shutdown(shutdown)
}

func env(k, fallback string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return fallback
}

var _ = net.IPv4len
