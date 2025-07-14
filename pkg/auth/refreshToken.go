package auth

import (
	"back-rex/pkg/utils"
	"errors"
	"net/http"

	"github.com/go-chi/render"
)

func refreshToken(w http.ResponseWriter, r *http.Request) {

	refresh_token, err := r.Cookie("refresh_cookie")
	if err != nil {
		render.Render(w, r, utils.ErrUnauthorizedRequest(err))
		return
	}

	refreshClaims, err := VerifyJWT(refresh_token.Value)
	if err != nil {
		render.Render(w, r, utils.ErrUnauthorizedRequest(err))
		return
	}

	currentToken := store.RechercherParSession(refreshClaims.session)
	if currentToken == nil {
		render.Render(w, r, utils.ErrUnauthorizedRequest(errors.New("session cookie inexistant")))
		return
	}

	if currentToken.Version != refreshClaims.version {
		// typique vol du refresh token. Supprime tous les tokens pour un utilisateur
		store.SupprimerParSujet(refreshClaims.UserID)
		render.Render(w, r, utils.ErrUnauthorizedRequest(errors.New("erreur version")))
		return
	}

	newVersion := refreshClaims.version + 1
	// regenere des nouveaux tokens, et les renvoi
	accessCokie, refreshCookie, err := GenerateJWTCookies(refreshClaims.UserID, &refreshClaims.session, &newVersion)
	if err != nil {
		render.Render(w, r, utils.ErrUnauthorizedRequest(err))
		return
	}

	err = store.MajParSession(refreshClaims.session, refreshCookie.Expires, newVersion)
	if err != nil {
		render.Render(w, r, utils.ErrUnauthorizedRequest(err))
		return
	}

	http.SetCookie(w, accessCokie)
	http.SetCookie(w, refreshCookie)
	w.WriteHeader(http.StatusNoContent)

}
