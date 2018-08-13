package routers

import (
	"DemoSecuringGOLangAPIWithJWT/controllers"
	"DemoSecuringGOLangAPIWithJWT/core/authentication"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

func SetHelloRoutes(router *mux.Router) *mux.Router {
	router.Handle("test/hello",
		negroni.New(
			negroni.HandlerFunc(authentication.RequireTokenAuthentication),
			negroni.HandlerFunc(controllers.HelloController),
		)).Methods("GET")

	return router
}
