package auth

import (
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"
)

var store TokenStore = TokenStore{}

func init() {
	go verifExpiration()
}

func verifExpiration() {
	// Crée un nouveau ticker qui émettra un signal toutes les 5 minutes
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	// Boucle infinie pour écouter les signaux du ticker
	for range ticker.C {
		store.SupprimerTokensExpires()
	}

}

type TokenJWT struct {
	Session        uuid.UUID
	Version        int
	DateExpiration time.Time
	Sujet          int
}

type TokenStore struct {
	mu     sync.RWMutex
	tokens []*TokenJWT
}

func (ts *TokenStore) AjouterToken(token TokenJWT) {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	t := token
	ts.tokens = append(ts.tokens, &t)
	sort.Slice(ts.tokens, func(i, j int) bool {
		return ts.tokens[i].Sujet < ts.tokens[j].Sujet
	})
}

func (ts *TokenStore) RechercherParSession(session uuid.UUID) *TokenJWT {
	ts.mu.RLock()
	defer ts.mu.RUnlock()

	for _, token := range ts.tokens {
		if token.Session == session {
			t := *token
			return &t
		}
	}
	return nil
}

func (ts *TokenStore) MajParSession(session uuid.UUID, exp time.Time, version int) error {
	fmt.Println("MajParSession")
	ts.mu.Lock()
	defer ts.mu.Unlock()

	for _, token := range ts.tokens {
		if token.Session == session {
			token.DateExpiration = exp
			token.Version = version
			return nil
		}
	}
	return errors.New("session inexistante")
}

func (ts *TokenStore) SupprimerParSujet(sujet int) {
	fmt.Println("SupprimerParSujet")
	ts.mu.Lock()
	defer ts.mu.Unlock()

	newTokens := []*TokenJWT{}
	for _, token := range ts.tokens {
		if token.Sujet != sujet {
			newTokens = append(newTokens, token)
		}
	}
	ts.tokens = newTokens
}

func (ts *TokenStore) SupprimerTokensExpires() {

	ts.mu.Lock()
	defer ts.mu.Unlock()

	now := time.Now()
	newTokens := []*TokenJWT{}
	fmt.Println("date: ", time.Now().String())
	for _, token := range ts.tokens {
		fmt.Println("\ttime: ", token.DateExpiration.String(), token.Version, token.Session, " time: ")
		if token.DateExpiration.After(now) {
			newTokens = append(newTokens, token)
		} else {
			fmt.Println("\t\texpiration de ", token.DateExpiration.String(), token.Version, token.Session)
		}
	}
	fmt.Println("supprimerTokensExpires: ", len(ts.tokens), len(newTokens))
	ts.tokens = newTokens
}
