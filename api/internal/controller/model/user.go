package model

import "github.com/hxcuber/friends-management/api/internal/repository/orm"

type User struct {
	UserID    int64  `boil:"user_id" json:"user_id" toml:"user_id" yaml:"user_id"`
	UserEmail string `boil:"user_email" json:"user_email" toml:"user_email" yaml:"user_email"`
}

type Users []User

func (u User) Orm() *orm.User {
	return &orm.User{
		UserID:    u.UserID,
		UserEmail: u.UserEmail,
	}
}
