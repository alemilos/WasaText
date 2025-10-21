package constants

const (
	MAX_UPLOAD_SIZE    = 2 * 1024 * 1024 // 2Mb
	PROFILE_PHOTO_PATH = "/profiles"
	GROUP_PHOTO_PATH   = "/groups"
	MESSAGE_PHOTO_PATH = "/messages"
)

// conversation types
const (
	CONV_GROUP   = "group"
	CONV_PRIVATE = "private"
)

// roles
const (
	ROLE_ADMIN  = "admin"
	ROLE_MEMBER = "member"
)
