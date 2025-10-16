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
	// Conversations
	rt.router.POST("/conversations", rt.wrap(rt.createConversation))
	rt.router.GET("/conversations", rt.wrap(rt.getMyConversations))
	rt.router.GET("/conversations/:id", rt.wrap(rt.getConversation))
	rt.router.POST("/conversations/:id/messages", rt.wrap(rt.sendMessage))
	rt.router.POST("/conversations/:id/read", rt.wrap(rt.readMessages))
	// Messages
	rt.router.POST("/messages/:id/forward", rt.wrap(rt.forwardMessage))
	rt.router.DELETE("/messages/:id", rt.wrap(rt.deleteMessage))
	// Comments
	rt.router.PUT("/messages/:id/comment", rt.wrap(rt.commentMessage))
	rt.router.DELETE("/messages/:id/comment", rt.wrap(rt.uncommentMessage))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
