package types

import "time"

type UserStore interface {
	ListUsers() ([]User, error)
	CreateUser(User) error
	GetUserById(int) (*User, error)
	GetUsersByEmail(string) ([]User, error)
	GetUsersByName(string) ([]User, error)
	UpdateUser(int, User) error
	DeleteUser(int) error
	GetUserTasks(int) ([]Task, error)
}

type TaskStore interface {
	ListTasks() ([]Task, error)
	CreateTask(Task) error
	GetTaskById(int) (*Task, error)
	GetTasksByQuery(string, string) ([]Task, error)
	UpdateTask(int, Task) error
	DeleteTask(int) error
}

type ProjectStore interface {
	ListProjects() ([]Project, error)
	CreateProject(Project) error
	GetProjectById(int) (*Project, error)
	GetProjectsByQuery(string, string) ([]Project, error)
	UpdateProject(int, Project) error
	DeleteProject(int) error
	GetProjectTasks(int) ([]Task, error)
}

type UserRole string

const (
	Admin     UserRole = "admin"
	Manager   UserRole = "manager"
	Developer UserRole = "developer"
)

type User struct {
	ID           int       `json:"id"`
	FullName     string    `json:"full_name"`
	Email        string    `json:"email"`
	RegisterDate time.Time `json:"register_date"`
	UserRole     UserRole  `json:"user_role"`
}

type TaskType string

const (
	Low    TaskType = "low"
	Medium TaskType = "medium"
	High   TaskType = "high"
)

type TaskPriority string

const (
	New        TaskPriority = "new"
	In_Process TaskPriority = "in_process"
	Done       TaskPriority = "done"
)

type Task struct {
	ID           int          `json:"id"`
	Title        string       `json:"title"`
	Descript     string       `json:"descript"`
	TaskType     TaskType     `json:"task_type"`
	TaskPriority TaskPriority `json:"task_priority"`
	UserId       int          `json:"user_id"`
	ProjectId    int          `json:"project_id"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
}

type Project struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Descript  string    `json:"descript"`
	ManagerId int       `json:"manager_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateUserPayload struct {
	FullName string   `json:"full_name" validate:"required"`
	Email    string   `json:"email" validate:"required"`
	UserRole UserRole `json:"user_role" validate:"required"`
}

type UpdateUserPayload struct {
	FullName string   `json:"full_name" validate:"omitempty"`
	UserRole UserRole `json:"user_role" validate:"omitempty"`
}

type CreateTaskPayload struct {
	Title        string       `json:"title" validate:"required"`
	Descript     string       `json:"descript" validate:"omitempty"`
	TaskType     TaskType     `json:"task_type" validate:"required"`
	TaskPriority TaskPriority `json:"task_priority" validate:"required"`
	UserId       int          `json:"user_id" validate:"required"`
	ProjectId    int          `json:"project_id" validate:"required"`
}

type UpdateTaskPayload struct {
	Title        string       `json:"title" validate:"required"`
	Descript     string       `json:"descript" validate:"required"`
	TaskType     TaskType     `json:"task_type" validate:"required"`
	TaskPriority TaskPriority `json:"task_priority" validate:"required"`
	UserId       int          `json:"user_id" validate:"required"`
	ProjectId    int          `json:"project_id" validate:"required"`
}

type CreateProjectPayload struct {
	Title     string `json:"title" validate:"required"`
	Descript  string `json:"descript" validate:"omitempty"`
	ManagerId int    `json:"manager_id" validate:"required"`
}

type UpdateProjectPayload struct {
	Title     string `json:"title" validate:"required"`
	Descript  string `json:"descript" validate:"required"`
	ManagerId int    `json:"manager_id" validate:"required"`
}
