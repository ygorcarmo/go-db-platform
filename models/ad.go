package models

import (
	"fmt"
	"log"

	"github.com/go-ldap/ldap/v3"
)

type LDAPConfig struct {
	ConnectionStr     string
	Username          string
	Password          string
	TopLevelDomain    string
	SecondLevelDomain string
}

// ldapsearch -H ldap://localhost:10389 -x -b "ou=people,dc=planetexpress,dc=com" -D "cn=admin,dc=planetexpress,dc=com" -w GoodNewsEveryone "(objectClass=inetOrgPerson)"

func (s LDAPConfig) Connect() (*ldap.Conn, error) {
	conn, err := ldap.DialURL(s.ConnectionStr)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	err = conn.Bind(fmt.Sprintf("cn=%s,dc=%s,dc=%s", s.Username, s.TopLevelDomain, s.SecondLevelDomain), s.Password)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to ldap")
	return conn, nil
}
