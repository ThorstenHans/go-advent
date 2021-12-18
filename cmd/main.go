package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/ThorstenHans/go-advent/pkg/app"
	"github.com/ThorstenHans/go-advent/pkg/automate"
)

func main() {
	p, ok := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT")
	if !ok {
		p = "8080"
	}
	a := app.New(fmt.Sprintf(":%s", p))
	a.Router.Use(app.ContentTypeJson)

	a.Router.HandleFunc("/CleanUpResourceGroups", automate.HandleTick).Methods(http.MethodPost)
	a.ListenAndServe()
}
