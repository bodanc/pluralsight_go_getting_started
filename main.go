package main

import (
	"gitlab.reynencourt.com/bogdan/go_learning/pluralsight/go_core_language/project/controllers"
	"net/http"
)

func main() {

	// we need to register our routing;
	controllers.RegisterControllers()
	// http.ListenAndServe() takes a serve multiplexer as its second parameter, an object which will handle all incoming
	// requests and the high level routing;
	http.ListenAndServe(":3000", nil)

}
