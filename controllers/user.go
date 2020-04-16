// we want to provide a way to route an http request to the userController;
// a handler (type) responds to an http request!
// polymorphism is the ability of an object to take many forms;

package controllers

import (
	"encoding/json"
	"gitlab.reynencourt.com/bogdan/go_learning/pluralsight/go_core_language/project/models"
	"net/http"
	"regexp"
	"strconv"
)

// an empty struct is useful if we have a common set of related behaviors that we want to associate together;
type someStruct struct {
}

// will handle 2 types of requests: resource requests on the users collection (all users) as well as allow us to
// manipulate individual users;
type userController struct {
	userIDPattern *regexp.Regexp
}

// this method signature requires a ResponseWriter object, as well as a Request object, from the http package;
// the method will work directly with the http information coming from the web request;
// because serveHTTP() receives a response writer as its first parameter, and a pointer to a request object as its
// second parameter, it automatically implements the Handler interface;
func (uc userController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// with this method, we are interacting with the outside world; we are writing back to the
	// http ResponseWriter object;
	if r.URL.Path == "/users" {
		switch r.Method {
		case http.MethodGet:
			// call into the getAll() method;
			uc.getAll(w, r)
		case http.MethodPost:
			// call into the post() method;
			uc.post(w, r)
		default:
			w.WriteHeader(http.StatusNotImplemented)
		}
	} else {
		// FindStringSubmatch() returns a slice of strings holding the text of the leftmost match of the regular
		// expression in r.URL.Path;
		// FindStringSubmatch() on the regular expression compares the incoming information (path) with ...
		// and will populate 'matches' with a slice containing all of the matches;
		matches := uc.userIDPattern.FindStringSubmatch(r.URL.Path)
		if len(matches) == 0 {
			w.WriteHeader(http.StatusNotFound)
		}

		id, err := strconv.Atoi(matches[1])
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
		}
		switch r.Method {
		case http.MethodGet:
			uc.get(id, w)
		case http.MethodPut:
			uc.put(id, w, r)
		case http.MethodDelete:
			uc.delete(id, w)
		default:
			w.WriteHeader(http.StatusNotImplemented)
		}
	}
}

// we declare a constructor function to set up an object for someone else to use;
// newUserController() will not return a value but rather a pointer to a userController object to avoid an unnecessary
// copy operation of the object;
func newUserController() *userController {
	// we are using the 'address of' operator on an object we're constructing directly and not declaring a variable;
	// this is a local variable; we're constructing it in the scope of the function and returing it's memory address;
	return &userController{
		userIDPattern: regexp.MustCompile(`^/users/(\d+)/?`),
	}
}

// getAll() will handle retrieving all the users from the []users collection in the models layer;
func (uc *userController) getAll(w http.ResponseWriter, r *http.Request) {
	encodeResponseAsJSON(models.GetUsers(), w)
}

// the get() method will accept the response writer from the ServeHTTP method;
func (uc *userController) get(id int, w http.ResponseWriter) {
	u, err := models.GetUserByID(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	encodeResponseAsJSON(u, w)
}

// add a new user to the []users collection in the models layer;
func (uc *userController) post(w http.ResponseWriter, r *http.Request) {
	u, err := uc.parseRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("could not parse the user object"))
		return
	}
	u, err = models.AddUser(u)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	encodeResponseAsJSON(u, w)
}

func (uc *userController) put(id int, w http.ResponseWriter, r *http.Request) {
	u, err := uc.parseRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("could not parse the user object"))
		return
	}

	if id != u.ID {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("the ID of the submitted user must match the ID in the URL"))
		return
	}

	u, err = models.UpdateUser(u)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	encodeResponseAsJSON(u, w)
}

func (uc *userController) delete(id int, w http.ResponseWriter) {
	err := models.RemoveUserById(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	w.WriteHeader(http.StatusOK)
}

// helper method which will convert the HTTP request into a JSON body and then convert it into a User object;
func (uc *userController) parseRequest(r *http.Request) (models.User, error) {
	bodyJson := json.NewDecoder(r.Body)
	var u models.User
	err := bodyJson.Decode(&u)
	if err != nil {
		return models.User{}, err
	}
	return u, nil
}
