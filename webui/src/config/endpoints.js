export const endpoints = {
	// Authentication
	doLogin: "/login",

	// Users
	setMyUsername: "/me/username",
	setMyPhoto: "/me/photo",
	getUsers: "/users",

	// Conversations
	createConversation: "/conversations",
	getMyConversations: "/conversations",
	getConversation: "/conversations/:id",
	sendMessage: "/conversations/:id/messages",
	readMessages: "/conversations/:id/read",

	// Messages
	forwardMessage: "/messages/:id/forward",
	deleteMessage: "/messages/:id",
	commentMessage: "/messages/:id/comment",
	uncommentMessage: "/messages/:id/comment",

	// Groups
	createGroup: "/groups",
	addToGroup: "/groups/:id/members",
	leaveGroup: "/groups/:id/members",
	setGroupName: "/groups/:id/name",
	setGroupPhoto: "/groups/:id/photo",
};
