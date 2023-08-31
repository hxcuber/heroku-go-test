package basic_success

type Response struct {
	Success        bool `json:"success"`
	HttpStatusCode int  `json:"-"`
}

func New(statusCode int) *Response {
	return &Response{
		Success:        true,
		HttpStatusCode: statusCode,
	}
}
