package users

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

func (s *Store) ListUsers() ([]types.User, error) {
	rows, err := s.db.Query("SELECT * FROM users")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	users := []types.User{}
	for rows.Next() {
		user, err := ScanRowIntoUser(rows)
		if err != nil {
			return nil, err
		}

		users = append(users, *user)
	}

	return users, nil
}

func (s *Store) CreateUser(user types.User) error {
	_, err := s.db.Exec("INSERT INTO users (fullName, email, userRole)"+
		"VALUES ($1, $2, $3)", user.FullName, user.Email, user.UserRole)

	if err != nil {
		return err
	}

	return nil
}

func (s *Store) GetUserById(userId int) (*types.User, error) {
	rows, err := s.db.Query("SELECT * FROM users WHERE id = $1", userId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	user := new(types.User)
	for rows.Next() {
		user, err = ScanRowIntoUser(rows)
		if err != nil {
			return user, nil
		}
	}

	if user.ID == 0 {
		return nil, fmt.Errorf("user not found")
	}

	return user, nil
}

func (s *Store) GetUsersByEmail(email string) ([]types.User, error) {
	rows, err := s.db.Query("SELECT * FROM users WHERE email ILIKE $1", "%"+email+"%")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	users := []types.User{}
	for rows.Next() {
		user, err := ScanRowIntoUser(rows)
		if err != nil {
			return nil, err
		}

		users = append(users, *user)
	}

	return users, nil
}

func (s *Store) GetUsersByName(name string) ([]types.User, error) {
	rows, err := s.db.Query("SELECT * FROM users WHERE fullName ILIKE $1", "%"+name+"%")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	users := []types.User{}
	for rows.Next() {
		user, err := ScanRowIntoUser(rows)
		if err != nil {
			return nil, err
		}

		users = append(users, *user)
	}

	return users, nil
}

func (s *Store) UpdateUser(userId int, user types.User) error {
	_, err := s.db.Exec("UPDATE users SET "+
		"fullName = $1, userRole = $2 WHERE id = $3", user.FullName, user.UserRole, userId)

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

func (s *Store) DeleteUser(userId int) error {
	_, err := s.db.Exec("DELETE FROM users WHERE id = $1", userId)

	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

func (s *Store) GetUserTasks(userId int) ([]types.Task, error) {
	rows, err := s.db.Query("SELECT * FROM tasks WHERE userId = $1", userId)

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

func ScanRowIntoUser(rows *sql.Rows) (*types.User, error) {
	user := new(types.User)

	err := rows.Scan(
		&user.ID,
		&user.FullName,
		&user.Email,
		&user.RegisterDate,
		&user.UserRole,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}
