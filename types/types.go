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

type UserRole string

const (
	Admin     UserRole = "admin"
	Manager   UserRole = "manager"
	Developer UserRole = "developer"
)

type User struct {
	ID           int       `json:"id"`
	FullName    string    `json:"full_name"`
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
	New        TaskType = "new"
	In_Process TaskType = "in_process"
	Done       TaskType = "done"
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

type CreateUserPayload struct {
	FullName string   `json:"full_name" validate:"required"`
	Email     string   `json:"email" validate:"required"`
	UserRole  UserRole `json:"user_role" validate:"required"`
}

type UpdateUserPayload struct {
	FullName string   `json:"full_name" validate:"required"`
	UserRole  UserRole `json:"user_role" validate:"required"`
}
