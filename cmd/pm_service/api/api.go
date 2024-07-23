package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/4lerman/pm_service/internal/service/projects"
	"github.com/4lerman/pm_service/internal/service/tasks"
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
	tasksRouter := subRouter.PathPrefix("/tasks").Subrouter()
	projectsRouter := subRouter.PathPrefix("/projects").Subrouter()

	usersStore := users.NewStore(s.db)
	usersService := users.NewHandler(usersStore)
	usersService.RegisterRoutes(usersRouter)

	tasksStore := tasks.NewStore(s.db)
	tasksService := tasks.NewHandler(tasksStore)
	tasksService.RegisterRoutes(tasksRouter)

	projectsStore := projects.NewStore(s.db)
	projectsService := projects.NewHandler(projectsStore)
	projectsService.RegisterRoutes(projectsRouter)

	log.Println("Listening on", s.addr)
	return http.ListenAndServe(s.addr, router)
}
