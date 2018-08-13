package main

import (
	"DemoSecuringGOLangAPIWithJWT/routers"
	"DemoSecuringGOLangAPIWithJWT/settings"
	"github.com/urfave/negroni"
	"net/http"
)

func main() {
	settings.Init()
	router := routers.InitRoutes()
	n := negroni.Classic()
	n.UseHandler(router)
	http.ListenAndServe(":5000", n)
}
