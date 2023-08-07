package request

type PutSubcriptionRequest struct {
	requestor string `json:"requestor"`
	target    string `json:"target"`
}
