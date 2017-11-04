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
	Route{"Catalog", "GET", "/services", ListService},
	Route{"Catalog", "GET", "/services/{service}/info", InfoService},
	Route{"Catalog", "POST", "/services/{service}/install", AddService},
	Route{"Catalog", "PATCH", "/services/{service}/start", StartService},
	Route{"Catalog", "PATCH", "/services/{service}/stop", StopService},
	Route{"Catalog", "PATCH", "/services/{service}/up", UpService},
	Route{"Catalog", "PATCH", "/services/{service}/down", DownService},
	Route{"Catalog", "PATCH", "/services/{service}/update", UpdateService},
	Route{"Catalog", "PATCH", "/services/{service}/backup", BackupService},

	Route{"Persons", "GET", "/persons", ldap.Persons},
	Route{"Persons", "POST", "/persons", ldap.AddPerson},
	Route{"Persons", "PUT", "/persons/{dn}", ldap.EditPerson},
	Route{"Persons", "PATCH", "/persons/{dn}/services/{serviceDn}", ldap.GrantService},
	Route{"Persons", "DELETE", "/persons/{dn}", ldap.DeletePerson},

	Route{"Groups", "GET", "/groups", ldap.Groups},
}
