package model

type Relationship struct {
	SenderID   int64  `boil:"sender_id" json:"sender_id" toml:"sender_id" yaml:"sender_id"`
	ReceiverID int64  `boil:"receiver_id" json:"receiver_id" toml:"receiver_id" yaml:"receiver_id"`
	Friends    bool   `boil:"friends" json:"friends" toml:"friends" yaml:"friends"`
	Status     string `boil:"status" json:"status" toml:"status" yaml:"status"`
}

type RelationshipSlice []Relationship
