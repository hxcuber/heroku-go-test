package relationship

import (
	"context"
	"github.com/hxcuber/friends-management/api/internal/controller/model"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func (i impl) UpdateRelationship(ctx context.Context, rela model.Relationship) (int64, error) {
	return rela.Orm().Update(ctx, i.dbConn, boil.Blacklist())
}
