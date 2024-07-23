package tasks

import (
	"database/sql"
	"fmt"

	"github.com/4lerman/pm_service/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) ListTasks() ([]types.Task, error) {
	rows, err := s.db.Query("SELECT * FROM tasks")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	tasks_list := []types.Task{}
	for rows.Next() {
		task, err := ScanRowIntoTask(rows)
		if err != nil {
			return nil, err
		}

		tasks_list = append(tasks_list, *task)
	}

	return tasks_list, nil
}

func (s *Store) CreateTask(task types.Task) error {
	_, err := s.db.Exec("INSERT INTO tasks (title, descript, taskType, taskPriority, userId, projectId)"+
		"VALUES ($1, $2, $3, $4, $5, $6)",
		task.Title, task.Descript, task.TaskType, task.TaskPriority, task.UserId, task.ProjectId)

	if err != nil {
		return err
	}

	return nil
}

func (s *Store) GetTaskById(taskId int) (*types.Task, error) {
	rows, err := s.db.Query("SELECT * FROM tasks WHERE id = $1", taskId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	task := new(types.Task)
	for rows.Next() {
		task, err = ScanRowIntoTask(rows)
		if err != nil {
			return task, nil
		}
	}

	if task.ID == 0 {
		return nil, fmt.Errorf("task not found")
	}

	return task, nil
}

func (s *Store) GetTasksByQuery(queryType string, query string) ([]types.Task, error) {
	var sqlQuery string

	switch queryType {
	case "title":
		sqlQuery = "SELECT * FROM tasks WHERE title ILIKE $1"
		query = "%"+query+"%"
	case "status":
		sqlQuery = "SELECT * FROM tasks WHERE taskPriority = $1"
	case "priority":
		sqlQuery = "SELECT * FROM tasks WHERE taskType = $1"
	case "assignee":
		sqlQuery = "SELECT * FROM tasks WHERE userId = $1"
	case "project":
		sqlQuery = "SELECT * FROM tasks WHERE projectId = $1"
	default:
		return nil, fmt.Errorf("invalid query type: %s", queryType)
	}

	rows, err := s.db.Query(sqlQuery, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	tasks_list := []types.Task{}
	for rows.Next() {
		task, err := ScanRowIntoTask(rows)
		if err != nil {
			return nil, err
		}

		tasks_list = append(tasks_list, *task)
	}

	return tasks_list, nil
}

func (s *Store) UpdateTask(taskId int, task types.Task) error {
	_, err := s.db.Exec("UPDATE tasks SET "+
		"title = $1, descript = $2, taskType = $3, taskPriority = $4, userId = $5, projectId = $6, updatedAt = NOW() "+
		"WHERE id = $7", 
		task.Title, task.Descript, task.TaskType, task.TaskPriority, task.UserId, task.ProjectId, taskId)

	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	return nil
}

func (s *Store) DeleteTask(taskId int) error {
	_, err := s.db.Exec("DELETE FROM tasks WHERE id = $1", taskId)

	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	return nil
}

func ScanRowIntoTask(rows *sql.Rows) (*types.Task, error) {
	task := new(types.Task)

	err := rows.Scan(
		&task.ID,
		&task.Title,
		&task.Descript,
		&task.TaskType,
		&task.TaskPriority,
		&task.UserId,
		&task.ProjectId,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return task, nil
}
