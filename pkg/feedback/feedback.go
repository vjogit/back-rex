package feedback

import (
	"back-rex/pkg/utils"
	"net/http"

	"github.com/go-chi/render"
)

// contient le crud.

func CreateFeedBack(w http.ResponseWriter, r *http.Request) {
	feedbackRequest := &FeedbackRequest{}

	err := render.Bind(r, feedbackRequest)
	if err != nil {
		render.Render(w, r, utils.ErrRender(err))
		return
	}

	// sauvegarde en bd.

	w.WriteHeader(http.StatusNoContent)

}
