package projects

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/4lerman/pm_service/types"
	"github.com/4lerman/pm_service/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type Handler struct {
	store types.ProjectStore
}

func NewHandler(store types.ProjectStore) *Handler {
	return &Handler{
		store: store,
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("", h.handleListProjects).Methods(http.MethodGet)
	router.HandleFunc("", h.handleCreateProject).Methods(http.MethodPost)
	router.HandleFunc("/search", h.handleProjectByQuery).Methods(http.MethodGet)
	router.HandleFunc("/{id}", h.handleGetProjectById).Methods(http.MethodGet)
	router.HandleFunc("/{id}", h.handleUpdateProject).Methods(http.MethodPut)
	router.HandleFunc("/{id}", h.handleDeleteProject).Methods(http.MethodDelete)
	router.HandleFunc("/{id}/tasks", h.handleGetProjectTasks).Methods(http.MethodGet)
}

func (h *Handler) handleListProjects(w http.ResponseWriter, r *http.Request) {
	projects, err := h.store.ListProjects()

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, projects)
}

func (h *Handler) handleCreateProject(w http.ResponseWriter, r *http.Request) {
	var payload types.CreateProjectPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		return
	}

	err := h.store.CreateProject(types.Project{
		Title:     payload.Title,
		Descript:  payload.Descript,
		ManagerId: payload.ManagerId,
	})

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]string{"msg": "Created successfully"})
}

func (h *Handler) handleProjectByQuery(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	var queryType, query string

	if title := queryParams.Get("title"); title != "" {
		queryType = "title"
		query = title
	} else if status := queryParams.Get("manager"); status != "" {
		queryType = "manager"
		query = status
	} else {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid query parameters"))
		return
	}

	projects, err := h.store.GetProjectsByQuery(queryType, query)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, projects)
}

func (h *Handler) handleGetProjectById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("id is not indicated"))
		return
	}

	projectId, _ := strconv.Atoi(id)

	project, err := h.store.GetProjectById(projectId)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("failed to get project by id: %v", err))
		return
	}

	utils.WriteJSON(w, http.StatusOK, project)
}

func (h *Handler) handleUpdateProject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("id is not indicated"))
		return
	}

	projectId, _ := strconv.Atoi(id)

	var payload types.UpdateProjectPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		return
	}

	err := h.store.UpdateProject(projectId, types.Project{
		Title:     payload.Title,
		Descript:  payload.Descript,
		ManagerId: payload.ManagerId,
	})

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"msg": "Updated successfully"})
}

func (h *Handler) handleDeleteProject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("id is not indicated"))
		return
	}

	projectId, _ := strconv.Atoi(id)

	if err := h.store.DeleteProject(projectId); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"msg": "Deleted successfully"})
}

func (h *Handler) handleGetProjectTasks(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("id is not indicated"))
		return
	}

	projectId, _ := strconv.Atoi(id)
	if _, err := h.store.GetProjectById(projectId); err != nil {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("failed to get project by id: %v", err))
		return
	}

	tasks_list, err := h.store.GetProjectTasks(projectId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, tasks_list)
}
