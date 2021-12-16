package postgres

import (
	"database/sql"

	"github.com/jmoiron/sqlx"

	pb "github.com/asilbek17071/ToDo/genproto"
)

type taskRepo struct {
	db *sqlx.DB
}
   	
func NewTaskRepo(db *sqlx.DB) *taskRepo {
	return &taskRepo{db: db}
}

func (r *taskRepo) Create(task pb.Task) (pb.Task, error) {
	var id int64
	err := r.db.QueryRow(`
        INSERT INTO tasks(assignee, title, summary, deadline, status)
        VALUES ($1,$2,$3,$4,$5) returning id`, task.Assignee, task.Title, task.Summary, task.Deadline, task.Status).Scan(&id)
	if err != nil {
		return pb.Task{}, err
	}

	task, err = r.Get(id)
	if err != nil {
		return pb.Task{}, err
	}

	return task, nil
}

func (r *taskRepo) Get(id int64) (pb.Task, error) {
	var task pb.Task
	err := r.db.QueryRow(`
        SELECT id, assignee, title, summary, deadline, status FROM tasks
        WHERE id=$1`, id).Scan(&task.Id, &task.Assignee, &task.Title, &task.Summary, &task.Deadline, &task.Status)
	if err != nil {
		return pb.Task{}, err
	}

	return task, nil
}

func (r *taskRepo) List(page, limit int64) ([]*pb.Task, int64, error) {
	offset := (page - 1) * limit
	// offset := (page - 1) * limit
	rows, err := r.db.Queryx(
		`SELECT id, assignee, title, summary, deadline, status FROM tasks  OFFSET $1  LIMIT $2`,
		offset, limit)
	if err != nil {
		return nil, 0, err
	}
	if err = rows.Err(); err != nil {
		return nil, 0, err
	}
	
	defer rows.Close()
	
	var (
		tasks []*pb.Task
		count int64
	)

	for rows.Next() {
		var task  pb.Task
		err = rows.Scan(&task.Id, &task.Assignee, &task.Title, &task.Summary, &task.Deadline, &task.Status)
		if err != nil {
			return nil, 0, err
		}
		tasks = append(tasks, &task)
	}

	err = r.db.QueryRow(`SELECT count(*) FROM tasks`).Scan(&count)
	if err != nil {
		return nil, 0, err
	}

	return tasks, count, nil

}

func (r *taskRepo) Update(task pb.Task) (pb.Task, error) {
	result, err := r.db.Exec(`UPDATE tasks SET assignee=$2, title=$3, summary=$4, deadline=$5, status=$6 WHERE id=$1`,
		task.Id, task.Assignee, task.Title, task.Summary, task.Deadline, task.Status)
	if err != nil {
		return pb.Task{}, err
	}

	if i, _ := result.RowsAffected(); i == 0 {
		return pb.Task{}, sql.ErrNoRows
	}

	task, err = r.Get(task.Id)
	if err != nil {
		return pb.Task{}, err
	}

	return task, nil
}

func (r *taskRepo) Delete(id int64) error {
	result, err := r.db.Exec(`DELETE FROM tasks WHERE id=$1`, id)
	if err != nil {
		return err
	}

	if i, _ := result.RowsAffected(); i == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *taskRepo) Overdue(t string, limit, page int64) ([]*pb.Task, int64, error) {
	offset := (page - 1) * limit
	rows, err := r.db.Queryx(
		`SELECT id, assignee, title, summary, deadline, status FROM tasks where deadline < $1  LIMIT $2 OFFSET $3`,
		t, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	if err = rows.Err(); err != nil {
		return nil, 0, err
	}
	defer rows.Close() // nolint:errcheck

	var (
		tasks []*pb.Task
		count int64
	)
	for rows.Next() {
		var task  pb.Task
		err = rows.Scan(&task.Id, &task.Assignee, &task.Title, &task.Summary, &task.Deadline, &task.Status)
		if err != nil {
			return nil, 0, err
		}
		tasks = append(tasks, &task)
	}

	err = r.db.QueryRow(`SELECT count(*) FROM tasks`).Scan(&count)
	if err != nil {
		return nil, 0, err
	}

	return tasks, count, nil
}
