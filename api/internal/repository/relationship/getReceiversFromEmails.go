package relationship

import (
	"context"
	"fmt"
	"github.com/hxcuber/friends-management/api/internal/controller/model"
	"github.com/hxcuber/friends-management/api/internal/repository/orm"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"strings"
)

func (i impl) GetReceiversFromEmails(ctx context.Context, senderEmail string, emails []string) (model.UserSlice, error) {
	sender, err := i.getUserByEmail(ctx, senderEmail)
	if err != nil {
		return nil, err
	}

	var inClauseBuilder strings.Builder
	inClauseBuilder.WriteString("(")
	for index, email := range emails {
		if index == 0 {
			inClauseBuilder.WriteString(fmt.Sprintf("'%s'", email))
		} else {
			inClauseBuilder.WriteString(fmt.Sprintf(",'%s'", email))
		}
	}
	inClauseBuilder.WriteString(")")

	query := fmt.Sprintf(
		"SELECT %s, %s FROM %s "+
			"WHERE %s IN $1 "+
			"EXCEPT "+
			"SELECT %s, %s FROM %s "+
			"JOIN %s ON %s=%s "+
			"WHERE %s=$2 "+
			"AND %s IN $1 "+
			"AND %s='%s'",
		orm.UserTableColumns.UserID,
		orm.UserTableColumns.UserEmail,
		orm.TableNames.Users,
		orm.UserTableColumns.UserEmail,
		orm.UserTableColumns.UserID,
		orm.UserTableColumns.UserEmail,
		orm.TableNames.Users,
		orm.TableNames.Relationships,
		orm.UserTableColumns.UserID,
		orm.RelationshipTableColumns.ReceiverID,
		orm.RelationshipTableColumns.SenderID,
		orm.UserTableColumns.UserEmail,
		orm.RelationshipTableColumns.Status,
		orm.SubscriptionStatusRBlockedS,
	)

	var finalUsers model.UserSlice
	err = orm.NewQuery(qm.SQL(query, inClauseBuilder.String(), sender.UserID)).Bind(ctx, i.dbConn, &finalUsers)
	if err != nil {
		return nil, err
	}

	return finalUsers, nil
}
