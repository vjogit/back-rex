package auth

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestStorage(t *testing.T) {
	store := TokenStore{}

	token1 := TokenJWT{
		Session:        uuid.New(),
		Version:        1,
		DateExpiration: time.Now().Add(-2 * time.Hour),
		Sujet:          2,
	}
	token2 := TokenJWT{
		Session:        uuid.New(),
		Version:        2,
		DateExpiration: time.Now().Add(2 * time.Hour),
		Sujet:          1,
	}

	session := uuid.New()

	dateToken3 := time.Now().Add(3 * time.Hour)
	token3 := TokenJWT{
		Session:        session,
		Version:        3,
		DateExpiration: dateToken3,
		Sujet:          3,
	}

	token4 := TokenJWT{
		Session:        uuid.New(),
		DateExpiration: time.Now().Add(-2 * time.Hour),
		Version:        4,
		Sujet:          3,
	}

	store.AjouterToken(token4)
	store.AjouterToken(token1)
	store.AjouterToken(token2)
	store.AjouterToken(token3)

	//  Suppression par sujet
	store.SupprimerParSujet(2)

	tokenTrouve2 := store.RechercherParSession(session)
	if tokenTrouve2 != nil {
		fmt.Printf("Token trouvé par session  %v: %v\n", session, tokenTrouve2)
	}

	// Suppression des tokens expirés
	store.SupprimerTokensExpires()

}
