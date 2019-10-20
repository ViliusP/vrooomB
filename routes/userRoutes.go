package routes

import "../users"

var userRoutes = Routes{
	Route{
		"Index",
		"GET",
		"/users/",
		users.GetUsers,
	},
	Route{
		"Index",
		"GET",
		"/users/{id}",
		users.GetUserByID,
	},
	Route{
		"Index",
		"DELETE",
		"/users/{id}",
		users.DeleteUserByID,
	},
}
