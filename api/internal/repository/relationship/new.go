package relationship

import "github.com/hxcuber/friends-management/api/pkg/db/pg"

type Repository interface {
}

type impl struct {
	dbConn pg.ContextExecutor
}

func New(dbConn pg.ContextExecutor) Repository {
	return impl{dbConn: dbConn}
}
