package restutil

import "net/http"

type ErrorResponse struct {
	Message string `json:"message"`
}

func (e *ErrorResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
