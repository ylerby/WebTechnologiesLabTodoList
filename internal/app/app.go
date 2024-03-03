package app

import (
	"context"
	"time"
)

func (a *Application) Run() error {
	err := a.server.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func (a *Application) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.server.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}
