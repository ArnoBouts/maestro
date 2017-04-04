package ldap

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	ldap "gopkg.in/ldap.v2"
)

type Group struct {
	Dn      string
	Cn      string
	Members []string
}

func Groups(writer http.ResponseWriter, request *http.Request) {

	groups, err := loadGroups()
	if err != nil {
		http.Error(writer, err.Error(), 500)
	}

	payload, err := json.Marshal(groups)
	if err != nil {
		http.Error(writer, err.Error(), 500)
	}

	writer.Header().Add("Content-Type", "application/json")
	writer.Write(payload)
}

func loadGroups() ([]Group, error) {
	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%s", os.Getenv("LDAP_HOST"), os.Getenv("LDAP_PORT")))
	if err != nil {
		return nil, err
	}
	defer l.Close()

	err = l.Bind(os.Getenv("LDAP_ADMIN_DN"), os.Getenv("LDAP_ADMIN_PASSWORD"))
	if err != nil {
		return nil, err
	}

	searchRequest := ldap.NewSearchRequest(
		"ou=groups,dc=home", // The base dn to search
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(&(objectClass=groupOfNames))", // The filter to apply
		[]string{"dn", "cn", "member"},  // A list attributes to retrieve
		nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		return nil, err
	}

	res := make([]Group, len(sr.Entries))
	for i, entry := range sr.Entries {
		res[i] = Group{entry.DN, entry.GetAttributeValue("cn"), entry.GetAttributeValues("member")}
	}

	return res, nil
}
