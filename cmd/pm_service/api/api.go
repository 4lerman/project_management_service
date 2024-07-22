package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/4lerman/pm_service/internal/service/users"
	"github.com/gorilla/mux"
)

type APIServer struct {
	addr string
	db   *sql.DB
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()
	subRouter := router.PathPrefix("/api/v1").Subrouter()

	usersRouter := subRouter.PathPrefix("/users").Subrouter()

	userStore := users.NewStore(s.db)
	userService := users.NewHandler(userStore)
	userService.RegisterRoutes(usersRouter)

	log.Println("Listening on", s.addr)
	return http.ListenAndServe(s.addr, router)
}
