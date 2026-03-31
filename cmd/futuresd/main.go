// Command futuresd starts the futures trading HTTP server.
package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/luxfi/futures/pkg/provider"
	"github.com/luxfi/futures/pkg/provider/envconfig"
	"github.com/luxfi/futures/pkg/types"
)

func main() {
	addr := os.Getenv("LISTEN_ADDR")
	if addr == "" {
		addr = ":8090"
	}

	reg := provider.NewRegistry()
	envconfig.Register(reg)

	slog.Info("futuresd starting", "providers", reg.List(), "addr", addr)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   corsOriginsFromEnv(),
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	r.Get("/providers", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string][]string{"providers": reg.List()})
	})

	r.Route("/v1", func(r chi.Router) {
		r.Get("/contracts/{provider}/{underlying}", contractsHandler(reg))
		r.Get("/quote/{provider}/{symbol}", quoteHandler(reg))
		r.Post("/order/{provider}/{accountID}", createOrderHandler(reg))
		r.Get("/positions/{provider}/{accountID}", positionsHandler(reg))
		r.Delete("/position/{provider}/{accountID}/{symbol}", closePositionHandler(reg))
		r.Get("/margin/{provider}/{accountID}/{symbol}", marginHandler(reg))
	})

	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)

	go func() {
		slog.Info("starting futuresd", "addr", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server failed", "err", err)
			os.Exit(1)
		}
	}()

	<-done
	slog.Info("shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("shutdown error", "err", err)
	}
}

func contractsHandler(reg *provider.Registry) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p, err := reg.Get(chi.URLParam(r, "provider"))
		if err != nil {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
			return
		}
		underlying := chi.URLParam(r, "underlying")
		contracts, err := p.GetFuturesContracts(r.Context(), underlying)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, contracts)
	}
}

func quoteHandler(reg *provider.Registry) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p, err := reg.Get(chi.URLParam(r, "provider"))
		if err != nil {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
			return
		}
		symbol := chi.URLParam(r, "symbol")
		quote, err := p.GetFuturesQuote(r.Context(), symbol)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, quote)
	}
}

func createOrderHandler(reg *provider.Registry) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p, err := reg.Get(chi.URLParam(r, "provider"))
		if err != nil {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
			return
		}
		accountID := chi.URLParam(r, "accountID")

		var req struct {
			Symbol      string `json:"symbol"`
			Side        string `json:"side"`
			Qty         string `json:"qty"`
			OrderType   string `json:"order_type"`
			LimitPrice  string `json:"limit_price,omitempty"`
			StopPrice   string `json:"stop_price,omitempty"`
			TimeInForce string `json:"time_in_force"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
			return
		}

		order, err := p.CreateFuturesOrder(r.Context(), accountID, &types.CreateFuturesOrderRequest{
			Symbol:      req.Symbol,
			Side:        req.Side,
			Qty:         req.Qty,
			OrderType:   req.OrderType,
			LimitPrice:  req.LimitPrice,
			StopPrice:   req.StopPrice,
			TimeInForce: req.TimeInForce,
		})
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusCreated, order)
	}
}

func positionsHandler(reg *provider.Registry) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p, err := reg.Get(chi.URLParam(r, "provider"))
		if err != nil {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
			return
		}
		accountID := chi.URLParam(r, "accountID")
		positions, err := p.GetFuturesPositions(r.Context(), accountID)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, positions)
	}
}

func closePositionHandler(reg *provider.Registry) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p, err := reg.Get(chi.URLParam(r, "provider"))
		if err != nil {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
			return
		}
		accountID := chi.URLParam(r, "accountID")
		symbol := chi.URLParam(r, "symbol")
		order, err := p.CloseFuturesPosition(r.Context(), accountID, symbol, nil)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, order)
	}
}

func marginHandler(reg *provider.Registry) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p, err := reg.Get(chi.URLParam(r, "provider"))
		if err != nil {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
			return
		}
		accountID := chi.URLParam(r, "accountID")
		symbol := chi.URLParam(r, "symbol")
		margin, err := p.GetFuturesMargin(r.Context(), accountID, symbol)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, margin)
	}
}

func corsOriginsFromEnv() []string {
	if v := os.Getenv("CORS_ALLOWED_ORIGINS"); v != "" {
		return strings.Split(v, ",")
	}
	return []string{
		"https://lux.exchange",
		"https://zoo.exchange",
		"https://pars.market",
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
