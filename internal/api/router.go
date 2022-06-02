package api

import (
	"github.com/gorilla/mux"
	// "github.com/prometheus/client_golang/prometheus"

	"go.pixelfactory.io/pkg/observability/log"

	"github.com/pixelfactoryio/crashlooper/internal/api/handlers"
	"github.com/pixelfactoryio/crashlooper/internal/api/middlewares"
)

// NewRouter returns a new mux.Router.
// It creates and register the metrics handler, the status handler and the default handler.
func NewRouter(logger log.Logger) *mux.Router {
	router := mux.NewRouter()
	router.Use(middlewares.Logging(logger))

	statusHandler := handlers.NewStatusHandler()
	router.PathPrefix("/checks/health").Handler(statusHandler)

	defaultHandler := handlers.NewDefaultHandler()
	router.PathPrefix("/").Handler(defaultHandler)

	return router
}
