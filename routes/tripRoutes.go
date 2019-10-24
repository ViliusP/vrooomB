package routes

import "../trips"

var tripRoutes = Routes{
	Route{
		"Get all trips by ID (LIMIT)",
		"GET",
		"/trips",
		trips.GetTrips,
		[]string{"limit", "{[0-9]+}", "offset", "{[0-9]+}"},
	},
	Route{
		"Get all users",
		"GET",
		"/trips",
		trips.GetTrips,
		[]string{},
	},
	Route{
		"Get all users trips",
		"GET",
		"/users/{id}/trips",
		trips.GetUserTrips,
		[]string{},
	},
	Route{
		"Delete trip by ID",
		"DELETE",
		"/trips/{id}",
		trips.DeleteTripByID,
		[]string{},
	},
	Route{
		"Update trip by ID",
		"PATCH",
		"/trips/{id}",
		trips.UpdateTripByID,
		[]string{},
	},
	Route{
		"Create trip by ID",
		"POST",
		"/trips",
		trips.CreateTrip,
		[]string{},
	},
}
