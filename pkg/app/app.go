package app

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type App struct {
	Address string
	Router  *mux.Router
}

func (a *App) ListenAndServe() error {
	if err := http.ListenAndServe(a.Address, a.Router); err != nil {
		log.Printf("HTTP server failed: %s", err)
		return err
	}
	return nil
}

func New(addr string) *App {
	return &App{
		Address: addr,
		Router:  mux.NewRouter(),
	}
}
