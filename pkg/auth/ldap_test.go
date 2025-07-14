package auth

import "testing"

func TestLdapPwd(t *testing.T) {
	pwd := "toto"
	LdapAuthenticate("ian.bertin@etu.mines-ales.fr", pwd)
	LdapAuthenticate("vlasak", pwd)
}
