package main

import (
	"net/http"
        "maestro/catalog"
        "maestro/ldap"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{"Catalog", "GET", "/catalog", catalog.List},

	Route{"Persons", "GET", "/persons", ldap.Persons},
	Route{"Persons", "POST", "/persons", ldap.AddPerson},
	Route{"Persons", "PUT", "/persons/{dn}", ldap.EditPerson},
	Route{"Persons", "PATCH", "/persons/{dn}/services/{serviceDn}", ldap.GrantService},
	Route{"Persons", "DELETE", "/persons/{dn}", ldap.DeletePerson},

	Route{"Groups", "GET", "/groups", ldap.Groups},
}
