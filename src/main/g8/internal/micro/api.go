package micro

import (
	"net/http"

	"github.com/Gympass/gcore/v3/glog"
	"github.com/Gympass/gcore/v3/middleware"
	"github.com/gorilla/handlers"
	"gopkg.in/DataDog/dd-trace-go.v1/contrib/gorilla/mux"
)

// Config for API v1
type Config struct {
	Logger     glog.Logger
	Router     *mux.Router
	Middleware middleware.GMiddlewareHandlerError
}

// NewAPI create API handler
func NewAPI(c Config) {
	demoHandler := NewHandler(NewService(NewRepository()), c.Logger)
	SetRoutes(demoHandler, c.Router, c.Middleware)
}

// SetRoutes for API handler
func SetRoutes(handler *Handler, router *mux.Router, mw middleware.GMiddlewareHandlerError) {
	r := router.PathPrefix("/v1").Subrouter()
	r.Handle(
		"/demo/{uid}",
		handlers.CompressHandler(
			mw.HandlerError(
				handler.Demo,
			),
		)).Methods(http.MethodGet)
}
