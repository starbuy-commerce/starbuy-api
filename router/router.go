package router

import (
	"authentication-service/router/routes"

	"github.com/gorilla/mux"
)

// Build vai retornar um router com as rotas configuradas
func Build() *mux.Router {
	return routes.Configure(mux.NewRouter())
}
