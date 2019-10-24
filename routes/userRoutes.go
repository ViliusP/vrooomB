package routes

import "../users"

var userRoutes = Routes{
	Route{
		"Get all users by ID",
		"GET",
		"/users",
		users.GetUsers,
		[]string{"limit", "{[0-9]+}", "offset", "{[0-9]+}"},
	},
	Route{
		"Get all users",
		"GET",
		"/users",
		users.GetUsers,
		[]string{},
	},
	Route{
		"Get user by ID",
		"GET",
		"/users/{id}",
		users.GetUserByID,
		[]string{},
	},
	Route{
		"Update user by ID",
		"PATCH",
		"/users/{id}",
		users.UpdateUserByID,
		[]string{},
	},
	Route{
		"Delete user by ID",
		"DELETE",
		"/users/{id}",
		users.DeleteUserByID,
		[]string{},
	},
	Route{
		"Create user",
		"POST",
		"/users",
		users.CreateUser,
		[]string{},
	},
}
