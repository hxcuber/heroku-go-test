package error

import "net/http"

func (resp *Response) Render(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(resp.HttpStatusCode)
	if resp.ErrMessage == "" {
		resp.ErrMessage = "no error message - unknown error\n"
	}
	return nil
}
