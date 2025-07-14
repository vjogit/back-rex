package auth

import (
	"back-rex/pkg/utils"
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/render"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type LogginError struct {
	Message string `json:"Message"`
}

type Credentials struct {
	Logging  string `json:"logging"`
	Password string `json:"password"`
}

func (a *Credentials) Bind(r *http.Request) error {
	return nil
}

func logging(w http.ResponseWriter, r *http.Request) {

	var err error

	credentials := &Credentials{}
	if err := render.Bind(r, credentials); err != nil {
		render.Render(w, r, utils.ErrRender(err))
		return
	}

	pgCtx := utils.GetPgCtx(r.Context())
	queries := New(pgCtx.Db)

	user, err := queries.GetUserByLogging(context.Background(), credentials.Logging)
	if err != nil {
		data := err.Error()
		if data == `no rows in result set` {
			badLogging := LogginError{
				Message: "BAD_PWD_OR_LOGGING",
			}
			render.Render(w, r, utils.ErrBadRequest(errors.New("erreur de validation"), badLogging))
			return

		}
		render.Render(w, r, utils.ErrRender(err))
		return

	}

	err = bcrypt.CompareHashAndPassword(user.PwdHash, []byte(credentials.Password))
	if err != nil {
		render.Render(w, r, utils.ErrRender(err))
		return
	}

	session := uuid.New()
	version := 1

	accessCokie, refreshCookie, err := GenerateJWTCookies(int(user.ID), &session, &version)
	if err != nil {
		render.Render(w, r, utils.ErrRender(err))
		return
	}

	http.SetCookie(w, accessCokie)
	http.SetCookie(w, refreshCookie)

	store.AjouterToken(TokenJWT{
		Session:        session,
		Version:        version,
		DateExpiration: refreshCookie.Expires,
		Sujet:          int(user.ID),
	})

	// obliger de convertir en GetUserByIdRow pour eviter de passe le pwd

	userResponse := GetUserByIdRow{
		ID:      user.ID,
		Logging: user.Logging,
		Nom:     user.Nom,
		Prenom:  user.Prenom,
		Roles:   user.Roles,
	}

	if err := render.Render(w, r, &userResponse); err != nil {
		render.Render(w, r, utils.ErrRender(err))
		return
	}

}
