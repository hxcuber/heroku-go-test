package router

import (
	"context"
	healthHandler "github.com/hxcuber/friends-management/api/internal/api/rest/health"
	relationshipHandler "github.com/hxcuber/friends-management/api/internal/api/rest/relationship"
	"github.com/hxcuber/friends-management/api/internal/controller/relationship"
	"github.com/hxcuber/friends-management/api/internal/controller/system"
)

// New creates and returns a new Router instance
func New(
	ctx context.Context,
	corsOrigin []string,
	isGQLIntrospectionOn bool,
	systemCtrl system.Controller,
	relationshipCtrl relationship.Controller,
) Router {
	return Router{
		ctx:                     ctx,
		corsOrigins:             corsOrigin,
		isGQLIntrospectionOn:    isGQLIntrospectionOn,
		healthRESTHandler:       healthHandler.New(systemCtrl),
		relationshipRESTHandler: relationshipHandler.New(relationshipCtrl),
	}
}
