package routes

import "../jwtauth"

var authRoutes = Routes{
	Route{
		"Sign in",
		"POST",
		"/signin",
		jwtauth.SignIn,
		[]string{},
	},
	Route{
		"Check JWT",
		"GET",
		"/check",
		jwtauth.CheckJWT,
		[]string{},
	},
}
