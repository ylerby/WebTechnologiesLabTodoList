package handlers

import (
	"net/http"

	"backend/internal/auth"
	"backend/internal/cache"
	"backend/internal/database"
	"backend/internal/middleware"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type Handler struct {
	database *database.Database
	cache    *cache.Cache
	logger   *zap.SugaredLogger
}

func New(database *database.Database, cache *cache.Cache, logger *zap.SugaredLogger) *Handler {
	return &Handler{
		database: database,
		cache:    cache,
		logger:   logger,
	}
}

func (h *Handler) InitRoutes() http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/", middleware.MethodValidationMiddleware(h.logger, http.MethodGet, h.Main))
	router.HandleFunc("/create_todo_list", middleware.MethodValidationMiddleware(h.logger, http.MethodPost,
		auth.AuthorizationMiddleware(h.CreateTodoList)))
	router.HandleFunc("/get_todo_lists", middleware.MethodValidationMiddleware(h.logger, http.MethodPost,
		auth.AuthorizationMiddleware(h.GetTodoListByTitle)))
	router.HandleFunc("/get_all_todo_list", middleware.MethodValidationMiddleware(h.logger, http.MethodGet,
		auth.AuthorizationMiddleware(h.GetAllTodoLists)))
	router.HandleFunc("/update_todo_list", middleware.MethodValidationMiddleware(h.logger, http.MethodPut,
		auth.AuthorizationMiddleware(h.UpdateTodoList)))
	router.HandleFunc("/delete_todo_list", middleware.MethodValidationMiddleware(h.logger, http.MethodDelete,
		auth.AuthorizationMiddleware(h.DeleteTodoList)))
	router.HandleFunc("/set_comment", middleware.MethodValidationMiddleware(h.logger, http.MethodPost,
		auth.AuthorizationMiddleware(h.SetComment)))
	router.HandleFunc("/login", middleware.MethodValidationMiddleware(h.logger, http.MethodPost, h.Login))
	router.HandleFunc("/register", middleware.MethodValidationMiddleware(h.logger, http.MethodPost, h.Register))
	return router
}
