package handlers

import (
	"backend/internal/auth"
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
	router.HandleFunc("/get_todo_lists", middleware.MethodValidationMiddleware(h.logger, http.MethodPost, h.GetTodoListByTitle))
	router.HandleFunc("/get_all_todo_lists", middleware.MethodValidationMiddleware(h.logger, http.MethodGet, h.GetAllTodoLists))
	router.HandleFunc("/update_todo_list", middleware.MethodValidationMiddleware(h.logger, http.MethodPut, h.UpdateTodoList))
	router.HandleFunc("/delete_todo_list", middleware.MethodValidationMiddleware(h.logger, http.MethodDelete, h.DeleteTodoList))
	router.HandleFunc("/protected", auth.AuthorizationMiddleware(h.ProtectedHandler))
	router.HandleFunc("/login", h.Login)
	return router
}
