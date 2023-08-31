package relationship

import (
	"context"
	"github.com/hxcuber/friends-management/api/internal/controller/model"
)

func (i impl) DeleteRelationship(ctx context.Context, rela model.Relationship) error {
	_, err := rela.Orm().Delete(ctx, i.dbConn)
	return err
}
