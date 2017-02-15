package ldap

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"

	ldap "gopkg.in/ldap.v2"
)

type Person struct {
	Dn        string `json:"dn"`
	Cn        string `json:"cn"`
	Sn        string `json:"sn"`
	Mail      string `json:"mail"`
	GidNumber int    `json:"gidNumber"`
	UidNumber int    `json:"uidNumber"`
}

func Persons(writer http.ResponseWriter, request *http.Request) {

	persons, err := loadPersons()
	if err != nil {
		http.Error(writer, err.Error(), 500)
	}

	payload, err := json.Marshal(persons)
	if err != nil {
		http.Error(writer, err.Error(), 500)
	}

	writer.Header().Add("Content-Type", "application/json")
	writer.Write(payload)
}

func AddPerson(writer http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)
	defer request.Body.Close()

	var p Person

	err := decoder.Decode(&p)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}

	err = addPerson(p)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
	writer.WriteHeader(http.StatusCreated)
	writer.Header().Add("Location", fmt.Sprintf("/persons/%s", p.Dn))
}

func EditPerson(writer http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)
	defer request.Body.Close()

	var p Person

	err := decoder.Decode(&p)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}

	err = editPerson(mux.Vars(request)["dn"], p)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
	writer.WriteHeader(http.StatusOK)
}

func GrantService(writer http.ResponseWriter, request *http.Request) {

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}

	allow, err := strconv.ParseBool(fmt.Sprintf("%s", body))
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}

	err = grantService(mux.Vars(request)["dn"], mux.Vars(request)["serviceDn"], allow)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
	writer.WriteHeader(http.StatusOK)
}

func DeletePerson(writer http.ResponseWriter, request *http.Request) {

	err := deletePerson(mux.Vars(request)["dn"])
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
	writer.WriteHeader(http.StatusOK)
}

func loadPersons() ([]Person, error) {
	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%s", os.Getenv("LDAP_HOST"), os.Getenv("LDAP_PORT")))
	if err != nil {
		return nil, err
	}
	defer l.Close()

	err = l.Bind("cn=admin,dc=home", "admin")
	if err != nil {
		return nil, err
	}

	searchRequest := ldap.NewSearchRequest(
		"ou=people,dc=home", // The base dn to search
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(&(objectClass=person))",                                    // The filter to apply
		[]string{"dn", "cn", "sn", "mail", "gidNumber", "uidNumber"}, // A list attributes to retrieve
		nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		return nil, err
	}

	res := make([]Person, len(sr.Entries))
	for i, entry := range sr.Entries {
		gidNumber, _ := strconv.Atoi(entry.GetAttributeValue("gidNumber"))
		uidNumber, _ := strconv.Atoi(entry.GetAttributeValue("uidNumber"))
		res[i] = Person{entry.DN, entry.GetAttributeValue("cn"), entry.GetAttributeValue("sn"), entry.GetAttributeValue("mail"), gidNumber, uidNumber}
	}

	return res, nil
}

func addPerson(p Person) error {
	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%s", os.Getenv("LDAP_HOST"), os.Getenv("LDAP_PORT")))
	if err != nil {
		return err
	}
	defer l.Close()

	err = l.Bind("cn=admin,dc=home", "admin")
	if err != nil {
		return err
	}

	req := ldap.NewAddRequest(p.Dn)
	req.Attribute("objectClass", []string{"top", "person", "organizationalPerson", "inetOrgPerson", "posixAccount"})
	req.Attribute("uid", []string{p.Cn})
	req.Attribute("homeDirectory", []string{""})
	req.Attribute("cn", []string{p.Cn})
	req.Attribute("sn", []string{p.Sn})
	req.Attribute("mail", []string{p.Mail})
	req.Attribute("gidNumber", []string{strconv.Itoa(p.GidNumber)})
	req.Attribute("uidNumber", []string{strconv.Itoa(p.UidNumber)})

	err = l.Add(req)
	if err != nil {
		return err
	}

	return nil
}

func deletePerson(dn string) error {
	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%s", os.Getenv("LDAP_HOST"), os.Getenv("LDAP_PORT")))
	if err != nil {
		return err
	}
	defer l.Close()

	err = l.Bind("cn=admin,dc=home", "admin")
	if err != nil {
		return err
	}

	req := ldap.NewDelRequest(dn, nil)
	err = l.Del(req)
	if err != nil {
		return err
	}

	return nil
}

func editPerson(dn string, p Person) error {

	if dn != p.Dn {
		return errors.New("It's not posible to modify Dn")
	}

	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%s", os.Getenv("LDAP_HOST"), os.Getenv("LDAP_PORT")))
	if err != nil {
		return err
	}
	defer l.Close()

	err = l.Bind("cn=admin,dc=home", "admin")
	if err != nil {
		return err
	}

	req := ldap.NewModifyRequest(p.Dn)
	req.Replace("uid", []string{p.Cn})
	req.Replace("homeDirectory", []string{""})
	req.Replace("cn", []string{p.Cn})
	req.Replace("sn", []string{p.Sn})
	req.Replace("mail", []string{p.Mail})
	req.Replace("gidNumber", []string{strconv.Itoa(p.GidNumber)})
	req.Replace("uidNumber", []string{strconv.Itoa(p.UidNumber)})

	err = l.Modify(req)
	if err != nil {
		return err
	}

	return nil
}

func grantService(dn string, serviceDn string, allow bool) error {
	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%s", os.Getenv("LDAP_HOST"), os.Getenv("LDAP_PORT")))
	if err != nil {
		return err
	}
	defer l.Close()

	err = l.Bind("cn=admin,dc=home", "admin")
	if err != nil {
		return err
	}

	req := ldap.NewModifyRequest(serviceDn)
	if allow {
		req.Add("member", []string{dn})
	} else {
		req.Delete("member", []string{dn})
	}

	err = l.Modify(req)
	if err != nil {
		return err
	}

	return nil
}
