package request

type GetRecipientsRequest struct {
	sender string `json:"sender"`
	text   string `json:"text"`
}
