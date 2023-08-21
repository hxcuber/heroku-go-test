package relationship

import (
	"context"
	"github.com/hxcuber/friends-management/api/internal/controller/model"
)

func (i impl) DeleteRelationship(ctx context.Context, rela model.Relationship) (int64, error) {
	return rela.Orm().Delete(ctx, i.dbConn)
}
