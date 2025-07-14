package auth

import (
	"fmt"
	"regexp"

	"github.com/go-ldap/ldap/v3"
)

// Vérifie l'authentification d'un utilisateur LDAP
func LdapAuthenticate(identifiant, password string) error {
	ldapURL := "ldap://ldap.mines-ales.fr:389"
	baseDN := "dc=ema,dc=fr"

	// Connexion au serveur LDAP
	l, err := ldap.DialURL(ldapURL)
	if err != nil {
		return fmt.Errorf("LDAP connection failed: %w", err)
	}
	defer l.Close()

	var filter string
	if IsValidEmail(identifiant) {
		filter = fmt.Sprintf("(mail=%s)", ldap.EscapeFilter(identifiant))
	} else {
		filter = fmt.Sprintf("(uid=%s)", ldap.EscapeFilter(identifiant))
	}

	searchRequest := ldap.NewSearchRequest(
		baseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 1, 0, false,
		filter,
		[]string{"*"}, // si nil, retourne ts les attibuts.
		nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		return fmt.Errorf("LDAP search failed: %w", err)
	}
	if len(sr.Entries) != 1 {
		return fmt.Errorf("Utilisateur non trouvé ou multiple")
	}
	userDN := sr.Entries[0].DN

	// Tente de se binder avec le DN et le mot de passe
	err = l.Bind(userDN, password)
	if err != nil {
		return fmt.Errorf("Authentification échouée: %w", err)
	}

	return nil // Authentification réussie
}

// IsValidEmail vérifie si identifiant est un mail valide
func IsValidEmail(identifiant string) bool {
	// Expression régulière simple pour valider un email
	var re = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(identifiant)
}
