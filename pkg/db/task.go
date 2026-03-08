package db

import (
	"fmt"
)

const defaultTasksLimit = 50

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTask(task *Task) (int64, error) {

	query := `
 INSERT INTO scheduler (date, title, comment, repeat)
 VALUES (?, ?, ?, ?)
 `

	res, err := DB.Exec(query,
		task.Date,
		task.Title,
		task.Comment,
		task.Repeat,
	)

	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

// возвращаем список задач, отсортированных по дате
func Tasks(limit int) ([]*Task, error) {
	if limit <= 0 {
		limit = defaultTasksLimit
	}

	query := fmt.Sprintf(`
 SELECT id, date, title, comment, repeat
 FROM scheduler
 ORDER BY date ASC
 LIMIT %d
 `, limit)

	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*Task

	for rows.Next() {
		var t Task
		err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, &t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// если задач нет
	if tasks == nil {
		tasks = []*Task{}
	}

	return tasks, nil
}

// получаем задачу по id
func GetTask(id string) (*Task, error) {
	if id == "" {
		return nil, fmt.Errorf("identifier is not specified")
	}

	var t Task
	query := `
 SELECT id, date, title, comment, repeat
 FROM scheduler
 WHERE id = ?
 `
	err := DB.QueryRow(query, id).Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, fmt.Errorf("task not found")
		}
		return nil, err
	}

	return &t, nil
}

// обновляем запись в базе
func UpdateTask(task *Task) error {
	query := `
 UPDATE scheduler
 SET date = ?, title = ?, comment = ?, repeat = ?
 WHERE id = ?
 `

	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("invalid ID for updating the task")
	}

	return nil
}

// удаляем задачу
func DeleteTask(id string) error {
	res, err := DB.Exec(`DELETE FROM scheduler WHERE id = ?`, id)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return fmt.Errorf("task not found")
	}

	return nil
}

// обновляем только дату задачи
func UpdateDate(next string, id string) error {
	res, err := DB.Exec(`UPDATE scheduler SET date = ? WHERE id = ?`, next, id)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return fmt.Errorf("task not found")
	}

	return nil
}
