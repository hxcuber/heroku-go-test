package senderText

type Request struct {
	Sender string `json:"sender"`
	Text   string `json:"text"`
}
