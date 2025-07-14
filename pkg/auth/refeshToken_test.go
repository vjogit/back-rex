package auth

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func TestRefres(t *testing.T) {

	r := chi.NewRouter()
	userID := 127

	accessCokie, err := GenerateJWTCookie(userID, "access_token", time.Duration(-30)*time.Second, "/", nil, nil)
	if err != nil {
		t.Fatal(err)
		return
	}

	session := uuid.New()
	fmt.Println(session)
	version := 1
	refreshCookie, err := GenerateJWTCookie(userID, "refresh_cookie", time.Duration(3600)*time.Second, "/refresh", &session, &version)
	if err != nil {
		t.Fatal(err)
		return
	}

	store.AjouterToken(TokenJWT{
		Session:        session,
		Version:        version,
		DateExpiration: refreshCookie.Expires,
		Sujet:          userID,
	})

	r.Get("/refresh", func(w http.ResponseWriter, r *http.Request) {
		refreshToken(w, r)
	})

	ts := httptest.NewServer(r)
	defer ts.Close()

	resp, err := sendTestRequest(ts, "GET", "/refresh", accessCokie, refreshCookie)

	if err != nil {
		t.Fatal(err)
		return
	}

	fmt.Println(resp)

}

func sendTestRequest(ts *httptest.Server, method, path string, accessCokie *http.Cookie, refreshCookie *http.Cookie) (*http.Response, error) {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Cookie", accessCokie.String())
	req.Header.Add("Cookie", refreshCookie.String())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
