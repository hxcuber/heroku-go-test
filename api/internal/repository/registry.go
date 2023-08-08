package repository

import (
	"context"
	"github.com/hxcuber/friends-management/api/internal/repository/relationship"
	"github.com/hxcuber/friends-management/api/internal/repository/system"
	"time"

	"github.com/hxcuber/friends-management/api/pkg/db/pg"

	"github.com/cenkalti/backoff/v4"
	pkgerrors "github.com/pkg/errors"
)

// Registry is the registry of all the domain specific repositories and also provides transaction capabilities.
type Registry interface {
	// System returns the system repo
	System() system.Repository
	// Relationship returns the relationship repo
	Relationship() relationship.Repository
	// DoInTx wraps operations within a db tx
	DoInTx(ctx context.Context, txFunc func(ctx context.Context, txRepo Registry) error, overrideBackoffPolicy backoff.BackOff) error
}

// New returns a new instance of Registry
func New(dbConn pg.BeginnerExecutor) Registry {
	return impl{
		dbConn:       dbConn,
		system:       system.New(dbConn),
		relationship: relationship.New(dbConn),
	}
}

type impl struct {
	dbConn       pg.BeginnerExecutor // Only used to start DB txns
	tx           pg.ContextExecutor  // Only used to keep track if txn has already been started to prevent devs from accidentally creating nested txns
	system       system.Repository
	relationship relationship.Repository
}

// System returns the system repo
func (i impl) System() system.Repository {
	return i.system
}

func (i impl) Relationship() relationship.Repository {
	return i.relationship
}

// DoInTx wraps operations within a db tx
func (i impl) DoInTx(ctx context.Context, txFunc func(ctx context.Context, txRepo Registry) error, overrideBackoffPolicy backoff.BackOff) error {
	if i.tx != nil {
		return pkgerrors.WithStack(errNestedTx)
	}

	if overrideBackoffPolicy == nil {
		overrideBackoffPolicy = pg.ExponentialBackOff(3, time.Minute)
	}

	return pg.TxWithBackOff(ctx, overrideBackoffPolicy, i.dbConn, func(tx pg.ContextExecutor) error {
		newI := impl{
			tx:     tx,
			system: system.New(tx),
		}
		return txFunc(ctx, newI)
	})
}
