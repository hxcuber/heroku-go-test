package relationship

import (
	"context"
	"github.com/hxcuber/friends-management/api/internal/controller/model"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func (i impl) CreateRelationship(ctx context.Context, rela model.Relationship) error {
	return rela.Orm().Insert(ctx, i.dbConn, boil.Blacklist())
}
