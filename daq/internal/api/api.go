package api

import (
	"context"
	"daq/internal/config"
	"daq/internal/rec"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type HandlerCtx struct {
	Records *rec.RecordSet
}

func (ctx *HandlerCtx) rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	_, _ = w.Write([]byte("Pulse is up and running!"))
}

func (ctx *HandlerCtx) recordsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(ctx.Records)
	if err != nil {
		logrus.Errorf("Failed to encode records: %v", err)
	}
}

func StartHTTPServer(cfg *config.PulseConfig, records *rec.RecordSet, ctxt context.Context) error {
	logrus.Infof("Starting Pulse API on port %d", cfg.Http.Port)
	ctx := HandlerCtx{Records: records}
	mux := http.NewServeMux()
	mux.HandleFunc("/", ctx.rootHandler)
	mux.HandleFunc("/records", ctx.recordsHandler)

	srv := http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Http.Port),
		Handler: mux,
	}

	go func() {
		err := srv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logrus.Fatalf("Error occurred in HTTP server: %v", err)
		}
	}()

	<-ctxt.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := srv.Shutdown(shutdownCtx)
	if err != nil {
		logrus.Errorf("Error occurred while shutting down HTTP server: %v", err)
		return nil
	}

	logrus.Infof("HTTP server stopped")
	return nil
}
