package db

import (
	"database/sql"
	"fmt"
	"strconv"
)

type Task struct {
	ID      string `json:"id,omitempty"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat,omitempty"`
}

func AddTask(task *Task) (int64, error) {
	var id int64

	query := `
		INSERT INTO scheduler (date, title, comment, "repeat")
		VALUES (?, ?, ?, ?)
	`

	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err == nil {
		id, err = res.LastInsertId()
	}

	return id, err
}

func Tasks(limit int) ([]*Task, error) {
	query := `
		SELECT id, date, title, comment, "repeat"
		FROM scheduler
		ORDER BY date ASC
		LIMIT ?
	`

	rows, err := DB.Query(query, limit)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	tasks := make([]*Task, 0)

	for rows.Next() {
		var task Task
		var id int64

		err := rows.Scan(&id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}

		task.ID = strconv.FormatInt(id, 10)
		tasks = append(tasks, &task)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func GetTask(id string) (*Task, error) {
	query := `
		SELECT id, date, title, comment, "repeat"
		FROM scheduler
		WHERE id = ?
	`

	var task Task
	var dbID int64

	err := DB.QueryRow(query, id).Scan(&dbID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("task not found")
		}
		return nil, err
	}

	task.ID = strconv.FormatInt(dbID, 10)

	return &task, nil
}

func UpdateTask(task *Task) error {
	query := `
		UPDATE scheduler
		SET date = ?, title = ?, comment = ?, "repeat" = ?
		WHERE id = ?
	`

	taskID, err := strconv.ParseInt(task.ID, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid task ID: %w", err)
	}

	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, taskID)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("task not found")
	}

	return nil
}

func DeleteTask(id string) error {
	taskID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid task ID: %w", err)
	}

	query := `DELETE FROM scheduler WHERE id = ?`
	res, err := DB.Exec(query, taskID)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("task not found")
	}

	return nil
}

func UpdateDate(next string, id string) error {
	taskID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid task ID: %w", err)
	}

	query := `UPDATE scheduler SET date = ? WHERE id = ?`
	res, err := DB.Exec(query, next, taskID)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("task not found")
	}

	return nil
}
