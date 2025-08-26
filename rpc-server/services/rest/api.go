package rest

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/ethereum/go-ethereum/log"

	"github.com/the-web3/rpc-server/common/httputil"
	"github.com/the-web3/rpc-server/config"
	"github.com/the-web3/rpc-server/database"
	"github.com/the-web3/rpc-server/services/rest/routes"
	"github.com/the-web3/rpc-server/services/rest/service"
)

const (
	SupportAssetPath = "/api/v1/get_support_symbol_list"
	MarketPricePath  = "/api/v1/get_market_price"
)

type ApiConfig struct {
	HttpServer   config.ServerConfig
	MetricServer config.ServerConfig
}

type API struct {
	router  *chi.Mux
	apiSvr  *httputil.HTTPServer
	db      *database.DB
	stopped atomic.Bool
}

func NewApi(ctx context.Context, cfg *config.Config) (*API, error) {
	out := &API{}
	if err := out.initFromConfig(ctx, cfg); err != nil {
		return nil, errors.Join(err, out.Stop(ctx))
	}
	return out, nil
}

func (a *API) initFromConfig(ctx context.Context, cfg *config.Config) error {
	if err := a.initDB(ctx, cfg); err != nil {
		return fmt.Errorf("failed to init DB: %w", err)
	}
	a.initRouter(cfg.RestServer, cfg)
	if err := a.startServer(cfg.RestServer); err != nil {
		return fmt.Errorf("failed to start API server: %w", err)
	}
	return nil
}

func (a *API) initRouter(conf config.ServerConfig, cfg *config.Config) {
	v := new(service.Validator)

	svc := service.NewHandleSvc(v, a.db.SupportToken, a.db.TokenPrice, a.db.SupportTokenExchange)
	apiRouter := chi.NewRouter()
	h := routes.NewRoutes(apiRouter, svc)

	apiRouter.Use(middleware.Timeout(time.Second * 12))
	apiRouter.Use(middleware.Recoverer)

	apiRouter.Get(fmt.Sprintf(SupportAssetPath), h.GetSupportSymbolList)
	apiRouter.Get(fmt.Sprintf(MarketPricePath), h.GetMarketPrice)

	a.router = apiRouter
}

func (a *API) initDB(ctx context.Context, cfg *config.Config) error {
	initDb, err := database.NewDB(ctx, cfg.MasterDB)
	if err != nil {
		log.Error("failed to connect to slave database", "err", err)
		return err
	}
	a.db = initDb
	return nil
}

func (a *API) Start(ctx context.Context) error {
	return nil
}

func (a *API) Stop(ctx context.Context) error {
	var result error
	if a.apiSvr != nil {
		if err := a.apiSvr.Stop(ctx); err != nil {
			result = errors.Join(result, fmt.Errorf("failed to stop API server: %w", err))
		}
	}
	if a.db != nil {
		if err := a.db.Close(); err != nil {
			result = errors.Join(result, fmt.Errorf("failed to close DB: %w", err))
		}
	}
	a.stopped.Store(true)
	log.Info("API service shutdown complete")
	return result
}

func (a *API) startServer(serverConfig config.ServerConfig) error {
	log.Debug("API server listening...", "port", serverConfig.Port)
	addr := net.JoinHostPort(serverConfig.Host, strconv.Itoa(serverConfig.Port))
	srv, err := httputil.StarHttpServer(addr, a.router)
	if err != nil {
		return fmt.Errorf("failed to start API server: %w", err)
	}
	log.Info("API server started", "addr", srv.Addr().String())
	a.apiSvr = srv
	return nil
}

func (a *API) Stopped() bool {
	return a.stopped.Load()
}
