package projects

import (
	"database/sql"
	"fmt"

	"github.com/4lerman/pm_service/internal/service/tasks"
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

func (s *Store) ListProjects() ([]types.Project, error) {
	rows, err := s.db.Query("SELECT * FROM projects")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	projects := []types.Project{}
	for rows.Next() {
		project, err := ScanRowIntoProject(rows)
		if err != nil {
			return nil, err
		}

		projects = append(projects, *project)
	}

	return projects, nil
}

func (s *Store) CreateProject(project types.Project) error {
	_, err := s.db.Exec("INSERT INTO projects (title, descript, managerId) VALUES ($1, $2, $3)",
		project.Title, project.Descript, project.ManagerId)

	if err != nil {
		return err
	}

	return nil
}

func (s *Store) GetProjectById(projectId int) (*types.Project, error) {
	rows, err := s.db.Query("SELECT * FROM projects WHERE id = $1", projectId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	project := new(types.Project)
	for rows.Next() {
		project, err = ScanRowIntoProject(rows)
		if err != nil {
			return project, nil
		}
	}

	if project.ID == 0 {
		return nil, fmt.Errorf("project not found")
	}

	return project, nil
}

func (s *Store) GetProjectsByQuery(queryType string, query string) ([]types.Project, error) {
	var sqlQuery string

	switch queryType {
	case "title":
		sqlQuery = "SELECT * FROM projects WHERE title ILIKE $1"
		query = "%" + query + "%"
	case "manager":
		sqlQuery = "SELECT * FROM projects WHERE managerId = $1"
	default:
		return nil, fmt.Errorf("invalid query type: %s", queryType)
	}

	rows, err := s.db.Query(sqlQuery, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	projects := []types.Project{}
	for rows.Next() {
		project, err := ScanRowIntoProject(rows)
		if err != nil {
			return nil, err
		}

		projects = append(projects, *project)
	}

	return projects, nil
}

func (s *Store) UpdateProject(projectId int, project types.Project) error {
	_, err := s.db.Exec("UPDATE projects SET "+
		"title = $1, descript = $2, managerId = $3, updatedAt = NOW() "+
		"WHERE id = $4", project.Title, project.Descript, project.ManagerId, projectId)

	if err != nil {
		return fmt.Errorf("failed to update project: %w", err)
	}

	return nil
}

func (s *Store) DeleteProject(projectId int) error {
	_, err := s.db.Exec("DELETE FROM projects WHERE id = $1", )

	if err != nil {
		return fmt.Errorf("failed to delete project: %w", err)
	}

	return nil
}

func (s *Store) GetProjectTasks(projectId int) ([]types.Task, error) {
	rows, err := s.db.Query("SELECT * FROM tasks WHERE projectId = $1", projectId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	tasks_list := []types.Task{}
	for rows.Next() {
		task, err := tasks.ScanRowIntoTask(rows)
		if err != nil {
			return nil, err
		}

		tasks_list = append(tasks_list, *task)
	}

	return tasks_list, nil
}

func ScanRowIntoProject(rows *sql.Rows) (*types.Project, error) {
	project := new(types.Project)

	err := rows.Scan(
		&project.ID,
		&project.Title,
		&project.Descript,
		&project.CreatedAt,
		&project.UpdatedAt,
		&project.ManagerId,
	)

	if err != nil {
		return nil, err
	}

	return project, nil
}
