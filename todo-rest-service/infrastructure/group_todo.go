package infrastructure

import (
	"database/sql"
	"time"

	"github.com/hryze/kakeibo-app-api/todo-rest-service/domain/model"
)

type GroupTodoRepository struct {
	*MySQLHandler
}

func NewGroupTodoRepository(mysqlHandler *MySQLHandler) *GroupTodoRepository {
	return &GroupTodoRepository{mysqlHandler}
}

func (r *GroupTodoRepository) GetDailyImplementationGroupTodoList(date time.Time, groupID int) ([]model.GroupTodo, error) {
	query := `
        SELECT
            id,
            posted_date,
            updated_date,
            implementation_date,
            due_date,
            todo_content,
            complete_flag,
            user_id
        FROM
            group_todo_list
        WHERE
            group_id = ?
        AND
            implementation_date = ?`

	rows, err := r.MySQLHandler.conn.Queryx(query, groupID, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	implementationGroupTodoList := make([]model.GroupTodo, 0)
	for rows.Next() {
		var implementationGroupTodo model.GroupTodo
		if err := rows.StructScan(&implementationGroupTodo); err != nil {
			return nil, err
		}

		implementationGroupTodoList = append(implementationGroupTodoList, implementationGroupTodo)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return implementationGroupTodoList, nil
}

func (r *GroupTodoRepository) GetDailyDueGroupTodoList(date time.Time, groupID int) ([]model.GroupTodo, error) {
	query := `
        SELECT
            id,
            posted_date,
            updated_date,
            implementation_date,
            due_date,
            todo_content,
            complete_flag,
            user_id
        FROM
            group_todo_list
        WHERE
            group_id = ?
        AND
            due_date = ?`

	rows, err := r.MySQLHandler.conn.Queryx(query, groupID, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	dueGroupTodoList := make([]model.GroupTodo, 0)
	for rows.Next() {
		var groupDueTodo model.GroupTodo
		if err := rows.StructScan(&groupDueTodo); err != nil {
			return nil, err
		}

		dueGroupTodoList = append(dueGroupTodoList, groupDueTodo)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return dueGroupTodoList, nil
}

func (r *GroupTodoRepository) GetMonthlyImplementationGroupTodoList(firstDay time.Time, lastDay time.Time, groupID int) ([]model.GroupTodo, error) {
	query := `
        SELECT
            id,
            posted_date,
            updated_date,
            implementation_date,
            due_date,
            todo_content,
            complete_flag,
            user_id
        FROM
            group_todo_list
        WHERE
            group_id = ?
        AND
            implementation_date >= ?
        AND
            implementation_date <= ?
        ORDER BY
            implementation_date`

	rows, err := r.MySQLHandler.conn.Queryx(query, groupID, firstDay, lastDay)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	implementationGroupTodoList := make([]model.GroupTodo, 0)
	for rows.Next() {
		var implementationGroupTodo model.GroupTodo
		if err := rows.StructScan(&implementationGroupTodo); err != nil {
			return nil, err
		}

		implementationGroupTodoList = append(implementationGroupTodoList, implementationGroupTodo)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return implementationGroupTodoList, nil
}

func (r *GroupTodoRepository) GetMonthlyDueGroupTodoList(firstDay time.Time, lastDay time.Time, groupID int) ([]model.GroupTodo, error) {
	query := `
        SELECT
            id,
            posted_date,
            updated_date,
            implementation_date,
            due_date,
            todo_content,
            complete_flag,
            user_id
        FROM
            group_todo_list
        WHERE
            group_id = ?
        AND
            due_date >= ?
        AND
            due_date <= ?
        ORDER BY
            due_date`

	rows, err := r.MySQLHandler.conn.Queryx(query, groupID, firstDay, lastDay)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	dueGroupTodoList := make([]model.GroupTodo, 0)
	for rows.Next() {
		var groupDueTodo model.GroupTodo
		if err := rows.StructScan(&groupDueTodo); err != nil {
			return nil, err
		}

		dueGroupTodoList = append(dueGroupTodoList, groupDueTodo)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return dueGroupTodoList, nil
}

func (r *GroupTodoRepository) GetExpiredGroupTodoList(dueDate time.Time, groupID int) (*model.ExpiredGroupTodoList, error) {
	query := `
        SELECT
            id,
            posted_date,
            updated_date,
            implementation_date,
            due_date,
            todo_content,
            complete_flag,
            user_id
        FROM
            group_todo_list
        WHERE
            group_id = ?
        AND
            complete_flag = b'0'
        AND
            due_date <= ?
        ORDER BY
            due_date`

	rows, err := r.MySQLHandler.conn.Queryx(query, groupID, dueDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	expiredGroupTodoList := model.ExpiredGroupTodoList{
		ExpiredGroupTodoList: make([]model.GroupTodo, 0),
	}
	for rows.Next() {
		var expiredGroupTodo model.GroupTodo
		if err := rows.StructScan(&expiredGroupTodo); err != nil {
			return nil, err
		}

		expiredGroupTodoList.ExpiredGroupTodoList = append(expiredGroupTodoList.ExpiredGroupTodoList, expiredGroupTodo)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &expiredGroupTodoList, nil
}

func (r *GroupTodoRepository) GetGroupTodo(groupTodoId int) (*model.GroupTodo, error) {
	query := `
        SELECT
            id,
            posted_date,
            updated_date,
            implementation_date,
            due_date,
            todo_content,
            complete_flag,
            user_id
        FROM
            group_todo_list
        WHERE
            id = ?`

	var groupTodo model.GroupTodo
	if err := r.MySQLHandler.conn.QueryRowx(query, groupTodoId).StructScan(&groupTodo); err != nil {
		return nil, err
	}

	return &groupTodo, nil
}

func (r *GroupTodoRepository) PostGroupTodo(groupTodo *model.GroupTodo, userID string, groupID int) (sql.Result, error) {
	query := `
        INSERT INTO group_todo_list
            (implementation_date, due_date, todo_content, user_id, group_id)
        VALUES
            (?,?,?,?,?)`

	result, err := r.MySQLHandler.conn.Exec(query, groupTodo.ImplementationDate, groupTodo.DueDate, groupTodo.TodoContent, userID, groupID)

	return result, err
}

func (r *GroupTodoRepository) PutGroupTodo(groupTodo *model.GroupTodo, groupTodoID int) error {
	query := `
        UPDATE
            group_todo_list
        SET 
            implementation_date = ?,
            due_date = ?,
            todo_content = ?,
            complete_flag = ?
        WHERE
            id = ?`

	_, err := r.MySQLHandler.conn.Exec(query, groupTodo.ImplementationDate, groupTodo.DueDate, groupTodo.TodoContent, groupTodo.CompleteFlag, groupTodoID)

	return err
}

func (r *GroupTodoRepository) DeleteGroupTodo(groupTodoID int) error {
	query := `
        DELETE
        FROM
            group_todo_list
        WHERE
            id = ?`

	_, err := r.MySQLHandler.conn.Exec(query, groupTodoID)

	return err
}

func (r *GroupTodoRepository) SearchGroupTodoList(groupTodoSqlQuery string) ([]model.GroupTodo, error) {
	rows, err := r.MySQLHandler.conn.Queryx(groupTodoSqlQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var searchGroupTodoList []model.GroupTodo
	for rows.Next() {
		var searchGroupTodo model.GroupTodo
		if err := rows.StructScan(&searchGroupTodo); err != nil {
			return nil, err
		}

		searchGroupTodoList = append(searchGroupTodoList, searchGroupTodo)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return searchGroupTodoList, nil
}
