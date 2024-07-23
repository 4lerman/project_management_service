package tasks

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
	store types.TaskStore
}

func NewHandler(store types.TaskStore) *Handler {
	return &Handler{
		store: store,
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("", h.handleListTasks).Methods(http.MethodGet)
	router.HandleFunc("", h.handleCreateTask).Methods(http.MethodPost)
	router.HandleFunc("/search", h.handleGetTaskByQuery).Methods(http.MethodGet)
	router.HandleFunc("/{id}", h.handleGetTaskById).Methods(http.MethodGet)
	router.HandleFunc("/{id}", h.handleUpdateTask).Methods(http.MethodPut)
	router.HandleFunc("/{id}", h.handleDeleteTask).Methods(http.MethodDelete)
}

func (h *Handler) handleListTasks(w http.ResponseWriter, r *http.Request) {
	users, err := h.store.ListTasks()

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, users)
}

func (h *Handler) handleCreateTask(w http.ResponseWriter, r *http.Request) {
	var payload types.CreateTaskPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		return
	}

	err := h.store.CreateTask(types.Task{
		Title:        payload.Title,
		Descript:     payload.Descript,
		TaskType:     payload.TaskType,
		TaskPriority: payload.TaskPriority,
		UserId:       payload.UserId,
		ProjectId:    payload.ProjectId,
	})

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]string{"msg": "Created successfully"})

}

func (h *Handler) handleGetTaskById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("id is not indicated"))
		return
	}

	taskId, _ := strconv.Atoi(id)

	user, err := h.store.GetTaskById(taskId)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("failed to get task by id: %v", err))
		return
	}

	utils.WriteJSON(w, http.StatusOK, user)
}

func (h *Handler) handleUpdateTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("id is not indicated"))
		return
	}

	taskId, _ := strconv.Atoi(id)

	var payload types.UpdateTaskPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		return
	}

	err := h.store.UpdateTask(taskId, types.Task{
		Title:        payload.Title,
		Descript:     payload.Descript,
		TaskType:     payload.TaskType,
		TaskPriority: payload.TaskPriority,
		UserId:       payload.UserId,
		ProjectId:    payload.ProjectId,
	})

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"msg": "Updated successfully"})
}

func (h *Handler) handleDeleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("id is not indicated"))
		return
	}

	taskId, _ := strconv.Atoi(id)

	if err := h.store.DeleteTask(taskId); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"msg": "Deleted successfully"})
}

func (h *Handler) handleGetTaskByQuery(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	var queryType, query string

	if title := queryParams.Get("title"); title != "" {
		queryType = "title"
		query = title
	} else if status := queryParams.Get("status"); status != "" {
		queryType = "status"
		query = status
	} else if priority := queryParams.Get("priority"); priority != "" {
		queryType = "priority"
		query = priority
	} else if assignee := queryParams.Get("assignee"); assignee != "" {
		queryType = "assignee"
		query = assignee
	} else if project := queryParams.Get("project"); project != "" {
		queryType = "project"
		query = project
	} else {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid query parameters"))
		return
	}

	tasks_list, err := h.store.GetTasksByQuery(queryType, query)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, tasks_list)
}
