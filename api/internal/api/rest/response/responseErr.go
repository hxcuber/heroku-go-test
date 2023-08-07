package response

type ResponseError struct {
	ResponseStatus
	errMessage string `json:"error_message"`
}
