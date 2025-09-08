package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"tekton-assistant/pkg/extractor"
)

// HandlerFunc defines a generic HTTP handler function type
type HandlerFunc func(w http.ResponseWriter, r *http.Request)

// httpServer wraps http.Server and modular handlers
type httpServer struct {
	*http.Server
	httpServerEndpoint string
	log                *log.Logger
	handlers           map[string]HandlerFunc
}

// NewHTTPServer creates a new httpServer with modular handlers
func NewHTTPServer(endpoint string, log *log.Logger) *httpServer {
	h := &httpServer{
		httpServerEndpoint: endpoint,
		log:                log,
		handlers:           make(map[string]HandlerFunc),
	}

	h.registerHandlers()
	h.initServer()
	return h
}

// registerHandlers registers all HTTP endpoints
func (h *httpServer) registerHandlers() {
	h.handlers["/taskrun/diagnose"] = h.handleExtract
	// Add more endpoints here if needed
}

// initServer wires handlers, metrics, CORS, and creates http.Server
func (h *httpServer) initServer() {
	mux := http.NewServeMux()
	for path, handler := range h.handlers {
		// Wrap with Prometheus metrics and CORS
		// handler := promhttp.InstrumentHandlerDuration(server.MetricLatency, mux)
		// handler = promhttp.InstrumentHandlerCounter(server.RequestsCount, handler)
		// handler = cors.Default().Handler(handler)
		mux.HandleFunc(path, handler)
	}

	h.Server = &http.Server{
		Addr:         h.httpServerEndpoint,
		Handler:      mux,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
}

// handleExtract handles the /extract endpoint
func (h *httpServer) handleExtract(w http.ResponseWriter, r *http.Request) {
	taskrunID := r.URL.Query().Get("taskrun_id")
	namespace := r.URL.Query().Get("namespace")
	if taskrunID == "" || namespace == "" {
		http.Error(w, "missing taskrun_id or namespace", http.StatusBadRequest)
		return
	}

	h.log.Printf("Extract request received: taskrun_id=%s, namespace=%s", taskrunID, namespace)

	result, err := extractor.ExtractTaskRunContext(r.Context(), taskrunID, namespace)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to extract context: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		h.log.Printf("Failed to encode response: %v", err)
	}
}

// startListener starts the HTTP server with graceful shutdown
func (h *httpServer) startListener(wg *sync.WaitGroup) {
	h.log.Printf("HTTP server listening on %s", h.httpServerEndpoint)

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM)
		<-sigint

		if err := h.Shutdown(context.Background()); err != nil {
			h.log.Printf("HTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed)
		h.log.Printf("stopped http server")
	}()

	wg.Add(1)
	go func() {
		if err := h.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			h.log.Fatalf("HTTP server failed: %v", err)
		}
		<-idleConnsClosed
		wg.Done()
		h.log.Printf("http server shutdown")
	}()
}
