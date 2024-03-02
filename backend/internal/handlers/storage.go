package handlers

import (
	"backend/internal/cache"
	"backend/internal/middleware"
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

	router.HandleFunc("/", middleware.MethodValidationMiddleware(h.logger, http.MethodGet, h.Main))
	router.HandleFunc("/create_todo_list", middleware.MethodValidationMiddleware(h.logger, http.MethodPost, h.CreateTodoList))
	router.HandleFunc("/get_all_todo", middleware.MethodValidationMiddleware(h.logger, http.MethodGet, h.GetAllTodoLists))
	return router
}
