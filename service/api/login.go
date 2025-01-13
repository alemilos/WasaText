package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// login endpoint handler
func (rt *_router) login(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Implement your login logic here
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Login successful"))
}
