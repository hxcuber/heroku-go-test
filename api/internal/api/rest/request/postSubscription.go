package request

type PostSubcriptionRequest struct {
	requestor string `json:"requestor"`
	target    string `json:"target"`
}
