package http

import "net/http"

func (e *ErrorResponse) Render(w http.ResponseWriter, r *http.Request) error {
	w.Header().Add("Content-Type", "application/json")
	return nil
}
