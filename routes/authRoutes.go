package routes

import "../jwtauth"

var authRoutes = Routes{
	Route{
		"Sign in",
		"Post",
		"/signin",
		jwtauth.SignIn,
		[]string{},
	},
}
