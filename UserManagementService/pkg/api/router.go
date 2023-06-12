package api

import (
	"advancedProg2Final/UserManagementService/pkg/api/middleware"
	"advancedProg2Final/UserManagementService/pkg/service"
	"github.com/gorilla/mux"
	"net/http"
)

func NewRouter(s *service.UserService) *mux.Router {
	r := mux.NewRouter()

	// Middleware
	r.Use(middleware.Logging)

	// User routes
	userHandler := NewUserHandler(s)
	userRoutes := r.PathPrefix("/api/user").Subrouter()
	userRoutes.HandleFunc("/create", userHandler.CreateUser).Methods(http.MethodPost)
	userRoutes.HandleFunc("/{id}", userHandler.GetUser).Methods(http.MethodGet)
	userRoutes.HandleFunc("/{id}", userHandler.UpdateUser).Methods(http.MethodPut)
	userRoutes.HandleFunc("/{id}", userHandler.DeleteUser).Methods(http.MethodDelete)

	return r
}
