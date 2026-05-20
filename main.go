package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	httptrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/net/http"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// version is overridden at build time via -ldflags "-X main.version=<sha>".
var version = "dev"

func main() {
	tracer.Start(
		tracer.WithService(envOr("DD_SERVICE", "svc-hello")),
		tracer.WithEnv(envOr("DD_ENV", "dev")),
		tracer.WithServiceVersion(envOr("DD_VERSION", version)),
	)
	defer tracer.Stop()

	mux := httptrace.NewServeMux()
	mux.HandleFunc("/healthz", healthz)
	mux.HandleFunc("/hello", hello)

	addr := envOr("LISTEN_ADDR", ":8080")
	srv := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("svc-hello listening on %s (version=%s)", addr, version)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	log.Println("shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
}

func healthz(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

func hello(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "world"
	}
	resp := map[string]string{
		"message": "hello, " + name,
		"version": version,
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
