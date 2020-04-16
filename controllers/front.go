package controllers

import (
	"encoding/json"
	"io"
	"net/http"
)

func RegisterControllers() {

	uc := newUserController()
	// any type that has a serveHTTP() method associated to it and has the correct function signature will implement this
	// interface;
	// we want our requests to be handled on the '/users' route;
	// we're expecting a handler interface for the 2nd parameter;
	http.Handle("/users", *uc)
	http.Handle("/users/", *uc)

}

func encodeResponseAsJSON(data interface{}, w io.Writer) {
	// newEncoder() creates an encoder which encodes go objects into json representations;
	enc := json.NewEncoder(w)
	enc.Encode(data)
}
