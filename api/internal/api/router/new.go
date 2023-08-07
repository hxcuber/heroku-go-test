package router

import (
	"context"
	"github.com/hxcuber/friends-management/api/internal/api/rest/healthHandler"
	"github.com/hxcuber/friends-management/api/internal/controller/systemController"
)

// New creates and returns a new Router instance
func New(
	ctx context.Context,
	corsOrigin []string,
	isGQLIntrospectionOn bool,
	systemCtrl systemController.Controller,
) Router {
	return Router{
		ctx:                  ctx,
		corsOrigins:          corsOrigin,
		isGQLIntrospectionOn: isGQLIntrospectionOn,
		healthRESTHandler:    healthHandler.New(systemCtrl),
	}
}
