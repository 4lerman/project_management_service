package tasks

import (
	"database/sql"

	"github.com/4lerman/pm_service/types"
)

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
