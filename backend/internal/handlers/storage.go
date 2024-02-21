package handlers

import (
	"backend/internal/cache"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
)

type Handler struct {
	cache  *cache.Cache
	logger *zap.SugaredLogger
}

func New(cache *cache.Cache, logger *zap.SugaredLogger) *Handler {
	return &Handler{
		cache:  cache,
		logger: logger,
	}
}

func (h *Handler) InitRoutes() http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/", h.Main).Methods(http.MethodGet)

	return router
}
