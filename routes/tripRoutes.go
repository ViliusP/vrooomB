package routes

import "../comments"

var commentRoutes = Routes{
	Route{
		"Create comment",
		"POST",
		"/users/{id_USER}/trips/{id_TRIP}/comment",
		comments.InsertComment,
		[]string{},
	},
}
