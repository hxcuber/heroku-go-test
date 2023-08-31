package sender_text

type Request struct {
	Sender string `json:"sender"`
	Text   string `json:"text"`
}
