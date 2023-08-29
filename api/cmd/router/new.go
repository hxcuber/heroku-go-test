package router

import (
	"context"
	"github.com/hxcuber/friends-management/api/internal/controller/relationship"
	"github.com/hxcuber/friends-management/api/internal/controller/system"
	"github.com/hxcuber/friends-management/api/internal/controller/user"
	healthHandler "github.com/hxcuber/friends-management/api/internal/handler/health"
	relationshipHandler "github.com/hxcuber/friends-management/api/internal/handler/relationship"
	userHandler "github.com/hxcuber/friends-management/api/internal/handler/user"
)

// New creates and returns a new Router instance
func New(
	ctx context.Context,
	corsOrigin []string,
	isGQLIntrospectionOn bool,
	systemCtrl system.Controller,
	userCtrl user.Controller,
	relationshipCtrl relationship.Controller,
) Router {
	return Router{
		ctx:                     ctx,
		corsOrigins:             corsOrigin,
		isGQLIntrospectionOn:    isGQLIntrospectionOn,
		healthRESTHandler:       healthHandler.New(systemCtrl),
		userRESTHandler:         userHandler.New(userCtrl),
		relationshipRESTHandler: relationshipHandler.New(relationshipCtrl),
	}
}
