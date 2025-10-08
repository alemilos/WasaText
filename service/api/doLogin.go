package api

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
)

// login endpoint handler
func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Login successful"))
}
