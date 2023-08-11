package model

type Relationship struct {
	SenderID   int64  `boil:"sender_id,bind" json:"sender_id" toml:"sender_id" yaml:"sender_id"`
	ReceiverID int64  `boil:"receiver_id,bind" json:"receiver_id" toml:"receiver_id" yaml:"receiver_id"`
	Friends    bool   `boil:"friends,bind" json:"friends" toml:"friends" yaml:"friends"`
	Status     string `boil:"status,bind" json:"status" toml:"status" yaml:"status"`
}

type RelationshipSlice []Relationship
