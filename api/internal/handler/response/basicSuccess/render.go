package basicSuccess

import (
	"net/http"
)

func (resp *Response) Render(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(resp.HttpStatusCode)
	return nil
}
