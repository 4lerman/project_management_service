package users

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
	store types.UserStore
}

func NewHandler(store types.UserStore) *Handler {
	return &Handler{
		store: store,
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("", h.handleListUsers).Methods(http.MethodGet)
	router.HandleFunc("", h.handleCreateUser).Methods(http.MethodPost)
	router.HandleFunc("/search", h.handleUserByNameOrEmail).Methods(http.MethodGet)
	router.HandleFunc("/{id}", h.handleGetUserById).Methods(http.MethodGet)
	router.HandleFunc("/{id}", h.handleUpdateUser).Methods(http.MethodPut)
	router.HandleFunc("/{id}", h.handleDeleteUser).Methods(http.MethodDelete)
	router.HandleFunc("/{id}/tasks", h.handleGetUserTasks).Methods(http.MethodGet)
}

// @Summary List all users
// @Description Get a list of all users
// @Tags Users
// @Accept  json
// @Produce  json
// @Success 200 {array} types.User
// @Failure 500 {object} map[string]string
// @Router /users [get]
func (h *Handler) handleListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.store.ListUsers()

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, users)
}

// @Summary Create a new user
// @Description Create a new user with the given details
// @Tags Users
// @Accept  json
// @Produce  json
// @Param user body types.CreateUserPayload true "User details"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users [post]
func (h *Handler) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	var payload types.CreateUserPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		return
	}

	err := h.store.CreateUser(types.User{
		FullName: payload.FullName,
		Email:    payload.Email,
		UserRole: payload.UserRole,
	})

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]string{"msg": "Created successfully"})

}

// @Summary Get user by ID
// @Description Get a user by their ID
// @Tags Users
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Success 200 {object} types.User
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /users/{id} [get]
func (h *Handler) handleGetUserById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("id is not indicated"))
		return
	}

	userId, _ := strconv.Atoi(id)

	user, err := h.store.GetUserById(userId)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("failed to get user by id: %v", err))
		return
	}

	utils.WriteJSON(w, http.StatusOK, user)
}

// @Summary Update user details
// @Description Update user details by ID
// @Tags Users
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Param user body types.UpdateUserPayload true "User details"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/{id} [put]
func (h *Handler) handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("id is not indicated"))
		return
	}

	userId, _ := strconv.Atoi(id)

	var payload types.UpdateUserPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		return
	}

	err := h.store.UpdateUser(userId, types.User{
		FullName: payload.FullName,
		UserRole: payload.UserRole,
	})

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"msg": "Updated successfully"})
}

// @Summary Delete user by ID
// @Description Delete a user by their ID
// @Tags Users
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/{id} [delete]
func (h *Handler) handleDeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("id is not indicated"))
		return
	}

	userId, _ := strconv.Atoi(id)

	if err := h.store.DeleteUser(userId); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"msg": "Deleted successfully"})
}

// @Summary Get user tasks
// @Description Get all tasks assigned to a user
// @Tags Users
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Success 200 {array} types.Task
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/{id}/tasks [get]
func (h *Handler) handleGetUserTasks(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("id is not indicated"))
		return
	}

	userId, _ := strconv.Atoi(id)
	if _, err := h.store.GetUserById(userId); err != nil {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("failed to get user by id: %v", err))
		return
	}

	tasks_list, err := h.store.GetUserTasks(userId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, tasks_list)
}

// @Summary Search users by name or email
// @Description Search users by name or email
// @Tags Users
// @Accept  json
// @Produce  json
// @Param name query string false "User name"
// @Param email query string false "User email"
// @Success 200 {array} types.User
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/search [get]
func (h *Handler) handleUserByNameOrEmail(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	email := r.URL.Query().Get("email")

	var users []types.User
	var err error

	if name != "" {
		users, err = h.store.GetUsersByName(name)
	} else if email != "" {
		users, err = h.store.GetUsersByEmail(email)
	} else {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("either name or email query parameter is required"))
		return
	}

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, users)
}
