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

// @Summary List all projects
// @Description Get a list of all projects
// @Tags Projects
// @Accept  json
// @Produce  json
// @Success 200 {array} types.Project
// @Failure 500 {object} map[string]string
// @Router /projects [get]
func (h *Handler) handleListProjects(w http.ResponseWriter, r *http.Request) {
	projects, err := h.store.ListProjects()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, projects)
}

// @Summary Create a new project
// @Description Create a new project with the given details
// @Tags Projects
// @Accept  json
// @Produce  json
// @Param project body types.CreateProjectPayload true "Project details"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /projects [post]
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

// @Summary Search projects by query
// @Description Search projects by title or manager ID
// @Tags Projects
// @Accept  json
// @Produce  json
// @Param title query string false "Project title"
// @Param manager query string false "Manager ID"
// @Success 200 {array} types.Project
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /projects/search [get]
func (h *Handler) handleProjectByQuery(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	var queryType, query string

	if title := queryParams.Get("title"); title != "" {
		queryType = "title"
		query = title
	} else if manager := queryParams.Get("manager"); manager != "" {
		queryType = "manager"
		query = manager
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

// @Summary Get project by ID
// @Description Get a project by its ID
// @Tags Projects
// @Accept  json
// @Produce  json
// @Param id path int true "Project ID"
// @Success 200 {object} types.Project
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /projects/{id} [get]
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

// @Summary Update project details
// @Description Update project details by ID
// @Tags Projects
// @Accept  json
// @Produce  json
// @Param id path int true "Project ID"
// @Param project body types.UpdateProjectPayload true "Project details"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /projects/{id} [put]
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

// @Summary Delete project by ID
// @Description Delete a project by its ID
// @Tags Projects
// @Accept  json
// @Produce  json
// @Param id path int true "Project ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /projects/{id} [delete]
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

// @Summary Get tasks by project ID
// @Description Get tasks associated with a project by project ID
// @Tags Projects
// @Accept  json
// @Produce  json
// @Param id path int true "Project ID"
// @Success 200 {array} types.Task
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /projects/{id}/tasks [get]
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
