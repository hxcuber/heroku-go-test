package router

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/hxcuber/friends-management/api/internal/api/rest/health"
	"github.com/hxcuber/friends-management/api/internal/api/rest/relationship"
	"net/http"
)

// Router defines the routes & handlers of the app
type Router struct {
	ctx                     context.Context
	corsOrigins             []string
	isGQLIntrospectionOn    bool
	healthRESTHandler       health.Handler
	relationshipRESTHandler relationship.Handler
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

	r.Get("/friends", rtr.relationshipRESTHandler.GetFriends())
	r.Get("/commonFriends", rtr.relationshipRESTHandler.GetCommonFriends())
	r.Get("/notification-receivers", rtr.relationshipRESTHandler.GetReceivers())
	r.Get("/_/ready", rtr.healthRESTHandler.CheckReadiness())
	r.Post("/subscribe", rtr.relationshipRESTHandler.CreateSubscription())
	r.Put("/subscribe", rtr.relationshipRESTHandler.CreateSubscription())
	return r
}
