package api

import (
	"encoding/json"
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
	"github.com/julienschmidt/httprouter"
)

type User struct {
	UserID    int64   `json:"userId"`
	Username  string  `json:"username"`
	PhotoPath *string `json:"photo_path"`
}

// a serializer, converter from database user to api response user.
func databaseToApiUsers(dbUsers []database.User) []User {
	apiUsers := make([]User, len(dbUsers))
	for i, u := range dbUsers {
		apiUsers[i] = User{
			UserID:    u.ID,
			Username:  u.Username,
			PhotoPath: u.PhotoPath,
		}
	}
	return apiUsers
}

type getUsersResponse struct {
	Users []User `json:"users"`
}

// get all the users registered to the application
func (rt *_router) getUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// === Get the users ===
	usersDB, err := rt.db.GetUsers()
	if err != nil {
		ctx.Logger.WithError(err).Error("failed to retrieve users")
		usersDB = []database.User{} // just return an empty []
	}

	users := databaseToApiUsers(usersDB)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(getUsersResponse{
		Users: users,
	})
}
