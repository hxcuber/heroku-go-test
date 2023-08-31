package requestor_target

type Request struct {
	Requestor string `json:"requestor"`
	Target    string `json:"target"`
}
