package api

import (
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) uncommentMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// extract the message_id from the params
	// extract the comment_id from the params
	// Check the comment existance
	// Make sure the user owns the comment_id (Unauthorized)
	// Remove the comment from the message
}
