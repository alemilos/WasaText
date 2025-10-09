package api

import (
	"net/http"
	"strconv"
	"strings"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

// httpRouterHandler is the signature for functions that accepts a reqcontext.RequestContext in addition to those
// required by the httprouter package.
type httpRouterHandler func(http.ResponseWriter, *http.Request, httprouter.Params, reqcontext.RequestContext)

// wrap parses the request and adds a reqcontext.RequestContext instance related to the request.
func (rt *_router) wrap(fn httpRouterHandler) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		reqUUID, err := uuid.NewV4()
		if err != nil {
			rt.baseLogger.WithError(err).Error("can't generate a request UUID")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		var ctx = reqcontext.RequestContext{
			ReqUUID: reqUUID,
		}

		// === Create a request-specific logger ===
		ctx.Logger = rt.baseLogger.WithFields(logrus.Fields{
			"reqid":     ctx.ReqUUID.String(),
			"remote-ip": r.RemoteAddr,
		})

		// === Authorization check ===
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			ctx.Logger.Warn("Missing Authorization header")
			http.Error(w, "Missing Authorization", http.StatusUnauthorized)
			return
		}

		const bearerPrefix = "Bearer "
		if !strings.HasPrefix(authHeader, bearerPrefix) {
			ctx.Logger.Warn("Invalid Authorization header format")
			http.Error(w, "Invalid Authorization Format", http.StatusUnauthorized)
			return
		}

		userIDStr := strings.TrimPrefix(authHeader, bearerPrefix)
		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			ctx.Logger.WithError(err).Warn("Invalid user ID in Authorization header")
			http.Error(w, "Invalid User ID", http.StatusUnauthorized)
			return
		}

		// === Fetch user from DB ===
		user, err := rt.db.GetUserById(userID)
		if err != nil {
			ctx.Logger.WithError(err).Warnf("User not found (ID=%d)", userID)
			http.Error(w, "User Not Found", http.StatusUnauthorized)
			return
		}

		ctx.User = user // Add the user to the context

		fn(w, r, ps, ctx)
	}
}
