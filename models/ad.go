package models

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"time"

	goLdap "github.com/go-ldap/ldap/v3"
)

type LDAP struct {
	ConnectionStr     string
	Username          string
	Password          string
	TopLevelDomain    string
	SecondLevelDomain string
	BaseGroup         string
	BaseGroupOU       string
	AdminGroup        string
	AdminGroupOU      string
	IsDefault         bool
	TimeOutInSecs     int
}

// ldapsearch -H ldap://localhost:10389 -x -b "ou=people,dc=planetexpress,dc=com" -D "cn=admin,dc=planetexpress,dc=com" -w GoodNewsEveryone "(objectClass=inetOrgPerson)"

// func (s *LDAP) Connect() (*goLdap.Conn, error) {
// 	conn, err := goLdap.DialURL(s.ConnectionStr)
// 	if err != nil {
// 		return nil, err
// 	}

// 	err = conn.Bind(fmt.Sprintf("cn=%s,dc=%s,dc=%s", s.Username, s.TopLevelDomain, s.SecondLevelDomain), s.Password)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return conn, nil
// }

func (s *LDAP) Connect() (*goLdap.Conn, error) {
	// Create a context with the specified timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(s.TimeOutInSecs)*time.Second)
	defer cancel()

	// Channel to handle connection and errors
	connChan := make(chan *goLdap.Conn, 1)
	errChan := make(chan error, 1)

	go func() {
		conn, err := goLdap.DialURL(s.ConnectionStr)
		if err != nil {
			errChan <- err
			return
		}

		err = conn.Bind(
			fmt.Sprintf("cn=%s,dc=%s,dc=%s", s.Username, s.TopLevelDomain, s.SecondLevelDomain),
			s.Password,
		)
		if err != nil {
			errChan <- err
			return
		}

		connChan <- conn
	}()

	select {
	case conn := <-connChan:
		return conn, nil
	case err := <-errChan:
		return nil, err
	case <-ctx.Done():
		return nil, fmt.Errorf("LDAP connection timed out after %d seconds", s.TimeOutInSecs)
	}
}

func (s *LDAP) Authenticate(conn *goLdap.Conn, username, password string) error {

	baseDN := fmt.Sprintf("dc=%s,dc=%s", s.TopLevelDomain, s.SecondLevelDomain)

	searchRequest := goLdap.NewSearchRequest(
		baseDN,
		goLdap.ScopeWholeSubtree, goLdap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(&(objectClass=organizationalPerson)(uid=%s))", goLdap.EscapeFilter(username)),
		[]string{"dn"},
		nil,
	)

	sr, err := conn.Search(searchRequest)
	if err != nil {
		return err
	}

	if len(sr.Entries) != 1 {
		return errors.New("user does not exist or too many entries returned")
	}

	userdn := sr.Entries[0].DN
	fmt.Println(sr.Entries[0])

	// Bind as the user to verify their password
	err = conn.Bind(userdn, password)
	if err != nil {
		return err
	}

	// Connect back as application user
	err = conn.Bind(fmt.Sprintf("cn=%s,dc=%s,dc=%s", s.Username, s.TopLevelDomain, s.SecondLevelDomain), s.Password)
	if err != nil {
		return err
	}

	return nil
}

func (*LDAP) GetGroupMembers(conn *goLdap.Conn, groupName, ou string) ([]string, error) {
	users := []string{}

	groupSearchRequest := goLdap.NewSearchRequest(
		fmt.Sprintf("ou=%s,dc=planetexpress,dc=com", ou),
		// search without OU filter
		// "dc=planetexpress,dc=com",
		goLdap.ScopeWholeSubtree,
		goLdap.NeverDerefAliases,
		0, 0, false,
		fmt.Sprintf("(&(objectClass=groupOfNames)(cn=%s))", groupName),
		[]string{"member"},
		nil,
	)

	result, err := conn.Search(groupSearchRequest)
	if err != nil {
		return users, nil
	}

	for _, entry := range result.Entries {
		members := entry.GetAttributeValues("member")
		for _, memberDN := range members {
			re := regexp.MustCompile(`uid=([^,]+)`)
			matches := re.FindStringSubmatch(memberDN)
			if len(matches) > 1 {
				users = append(users, matches[1])
			}
		}
	}

	return users, nil
}
