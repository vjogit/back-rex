package feedback

import "net/http"

type FeedbackRequest struct {
	Message string `json:"message"`
}

// Bind implements render.Binder.
func (f *FeedbackRequest) Bind(r *http.Request) error {
	return nil
}
