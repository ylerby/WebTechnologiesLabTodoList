package app

import "net/http"

const (
	serverAddress = ":8080"
)

type Application struct {
	server *http.Server
}

func New(router http.Handler) *Application {
	return &Application{
		server: &http.Server{
			Addr:    serverAddress,
			Handler: router,
		},
	}
}
