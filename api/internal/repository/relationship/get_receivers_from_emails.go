package relationship

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/hxcuber/friends-management/api/internal/controller/model"
	"github.com/hxcuber/friends-management/api/internal/repository/orm"
	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"strings"
)

func (i impl) GetReceiversFromEmails(ctx context.Context, sender model.User, emails []string) (model.Users, error) {
	if len(emails) == 0 {
		return nil, nil
	}

	var inClauseBuilder strings.Builder
	inClauseBuilder.WriteString("(")
	var emailsInterface []interface{}
	for index, email := range emails {
		if index == 0 {
			inClauseBuilder.WriteString(fmt.Sprintf("'%s'", email))
		} else {
			inClauseBuilder.WriteString(fmt.Sprintf(",'%s'", email))
		}
		emailsInterface = append(emailsInterface, email)
	}
	inClauseBuilder.WriteString(")")

	/*
		SQL injection is avoided because the emails are verified at the controller stage,
		and the rest are all constants.
	*/
	query := fmt.Sprintf(
		"SELECT %s, %s FROM %s "+
			"WHERE %s IN %s "+
			"EXCEPT "+
			"SELECT %s, %s FROM %s "+
			"JOIN %s ON %s=%s "+
			"WHERE %s=%d "+
			"AND %s IN %s "+
			"AND %s='%s'",
		orm.UserTableColumns.UserID,
		orm.UserTableColumns.UserEmail,
		orm.TableNames.Users,
		orm.UserTableColumns.UserEmail,
		inClauseBuilder.String(),
		orm.UserTableColumns.UserID,
		orm.UserTableColumns.UserEmail,
		orm.TableNames.Users,
		orm.TableNames.Relationships,
		orm.UserTableColumns.UserID,
		orm.RelationshipTableColumns.ReceiverID,
		orm.RelationshipTableColumns.SenderID,
		sender.UserID,
		orm.UserTableColumns.UserEmail,
		inClauseBuilder.String(),
		orm.RelationshipTableColumns.Status,
		orm.StatusRBlockedS,
	)

	var finalUsers model.Users
	err := orm.NewQuery(qm.SQL(query)).Bind(ctx, i.dbConn, &finalUsers)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return finalUsers, nil
}
