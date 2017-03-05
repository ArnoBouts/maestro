package main

import (
	"maestro/catalog"
	"maestro/ldap"
	"net/http"
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
	Route{"Catalog", "PATCH", "/catalog/{service}/start", catalog.StartService},
	Route{"Catalog", "PATCH", "/catalog/{service}/stop", catalog.StopService},
	Route{"Catalog", "PATCH", "/catalog/{service}/up", catalog.UpService},
	Route{"Catalog", "PATCH", "/catalog/{service}/down", catalog.DownService},

	Route{"Persons", "GET", "/persons", ldap.Persons},
	Route{"Persons", "POST", "/persons", ldap.AddPerson},
	Route{"Persons", "PUT", "/persons/{dn}", ldap.EditPerson},
	Route{"Persons", "PATCH", "/persons/{dn}/services/{serviceDn}", ldap.GrantService},
	Route{"Persons", "DELETE", "/persons/{dn}", ldap.DeletePerson},

	Route{"Groups", "GET", "/groups", ldap.Groups},
}
