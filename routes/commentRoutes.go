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
	Route{
		"Delete comment",
		"DELETE",
		"/users/{id_USER}/trips/{id_TRIP}/comment/{id_COMMENT}",
		comments.DeleteCommentByID,
		[]string{},
	},
	Route{
		"Get comments",
		"GET",
		"/users/{id_USER}/trips/{id_TRIP}/comment",
		comments.GetComments,
		[]string{},
	},
}
