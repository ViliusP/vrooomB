package routes

import (
	"net/http"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
	Queries     []string
}

type Routes []Route

func composeAllRoutes() Routes {
	var allRoutes = Routes{}
	allRoutes = append(allRoutes, userRoutes...)
	// allRoutes = append(allRoutes, tripRoutes...)
	return allRoutes
}
