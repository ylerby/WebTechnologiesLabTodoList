package app

import "net/http"

type Application struct {
	server *http.Server
}

func New(router http.Handler) *Application {
	return &Application{
		server: &http.Server{
			Addr:    "localhost:8080",
			Handler: router,
		},
	}
}
