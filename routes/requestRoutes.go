package routes

import (
	"../requests"
)

var requestRoutes = Routes{
	Route{
		"Get all users requests",
		"GET",
		"/users/{id}/requests",
		requests.GetUserRequests,
		[]string{},
	},
	Route{
		"Get all users requests (LIMIT)",
		"GET",
		"/users/{id}/requests",
		requests.GetUserRequests,
		[]string{"limit", "{[0-9]+}", "offset", "{[0-9]+}"},
	},
	Route{
		"Get all user's trip requests",
		"GET",
		"/trips/{id}/requests",
		requests.GetTripRequests,
		[]string{},
	},
	Route{
		"Delete request by ID",
		"DELETE",
		"/requests/{id}",
		requests.DeleteRequestByID,
		[]string{},
	},
	Route{
		"Update request status by ID",
		"PATCH",
		"/request/{id}",
		requests.UpdateRequestByID,
		[]string{},
	},
	Route{
		"Get requests status types",
		"GET",
		"/requests/statuses",
		requests.GetStatuses,
		[]string{},
	},
}
