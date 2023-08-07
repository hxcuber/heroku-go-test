package router

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"

	"github.com/hxcuber/friends-management/api/internal/api/rest/healthHandler"
)

// Router defines the routes & handlers of the app
type Router struct {
	ctx                  context.Context
	corsOrigins          []string
	isGQLIntrospectionOn bool
	healthRESTHandler    healthHandler.Handler
}

// Handler returns the Handler for use by the server
func (rtr Router) Handler() http.Handler {
	r := chi.NewRouter()
	// TODO: add middleware here
	r.Use(
		// render.SetContentType(render.ContentTypeJSON), // set content-type headers as application/json
		middleware.Logger, // log relationship request calls
		// middleware.DefaultCompress, // compress results, mostly gzipping assets and json
		middleware.StripSlashes, // match paths with a trailing slash, strip it, and continue routing through the mux
		middleware.Recoverer,    // recover from panics without crashing server
	)

	r.Get("/_/ready", rtr.healthRESTHandler.CheckReadiness())
	return r
}
