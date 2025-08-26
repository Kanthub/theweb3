package routes

import (
	"github.com/go-chi/chi/v5"

	"github.com/the-web3/rpc-server/services/rest/service"
)

type Routes struct {
	router *chi.Mux
	srv    service.RestService
}

func NewRoutes(r *chi.Mux, srv service.RestService) Routes {
	return Routes{
		router: r,
		srv:    srv,
	}
}
