package model

import "github.com/hxcuber/friends-management/api/internal/repository/orm"

type Relationship struct {
	ReceiverID int64  `boil:"receiver_id" json:"receiver_id" toml:"receiver_id" yaml:"receiver_id"`
	SenderID   int64  `boil:"sender_id" json:"sender_id" toml:"sender_id" yaml:"sender_id"`
	Status     string `boil:"status" json:"status" toml:"status" yaml:"status"`
}

type Relationships []Relationship

func (r Relationship) Orm() *orm.Relationship {
	return &orm.Relationship{
		ReceiverID: r.ReceiverID,
		SenderID:   r.SenderID,
		Status:     r.Status,
	}
}
