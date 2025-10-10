package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.POST("/login", rt.doLogin)
	// Users
	rt.router.PUT("/me/username", rt.wrap(rt.setMyUserName))
	rt.router.POST("/me/photo", rt.wrap(rt.setMyPhoto))
	rt.router.GET("/users", rt.wrap(rt.getUsers))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
