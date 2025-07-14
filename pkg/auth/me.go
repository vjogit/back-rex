package auth

import (
	"back-rex/pkg/utils"
	"context"
	"net/http"

	"github.com/go-chi/render"
)

func GetUserFromCookie(r *http.Request) (*GetUserByIdRow, error) {
	access_token, err := r.Cookie("access_token")
	if err != nil {
		return nil, err
	}

	accessClaims, err := VerifyJWT(access_token.Value)
	if err != nil {
		return nil, err
	}

	pgCtx := utils.GetPgCtx(r.Context())
	queries := New(pgCtx.Db)

	user, err := queries.GetUserById(context.Background(), int32(accessClaims.UserID))
	if err != nil {
		return nil, err
	}

	return &user, nil

}

func me(w http.ResponseWriter, r *http.Request) {
	user, err := GetUserFromCookie(r)

	if err != nil {
		render.Render(w, r, utils.ErrUnauthorizedRequest(err))
		return
	}

	if err := render.Render(w, r, user); err != nil {
		render.Render(w, r, utils.ErrRender(err))
		return
	}

}
