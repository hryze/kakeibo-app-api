package model

import (
	"database/sql/driver"
	"errors"
	"time"
)

type TodoList struct {
	ImplementationTodoList []Todo `json:"implementation_todo_list,omitempty"`
	DueTodoList            []Todo `json:"due_todo_list,omitempty"`
}

type Todo struct {
	ID                 int       `json:"id"                  db:"id"`
	PostedDate         time.Time `json:"posted_date"         db:"posted_date"`
	ImplementationDate Date      `json:"implementation_date" db:"implementation_date"`
	DueDate            Date      `json:"due_date"            db:"due_date"`
	TodoContent        string    `json:"todo_content"        db:"todo_content"`
	CompleteFlag       BitBool   `json:"complete_flag"       db:"complete_flag"`
}

type Date struct {
	time.Time
}

type BitBool bool

func NewTodoList(implementationTodoListTodoList []Todo, dueTodoList []Todo) TodoList {
	return TodoList{
		ImplementationTodoList: implementationTodoListTodoList,
		DueTodoList:            dueTodoList,
	}
}

func (d *Date) Scan(value interface{}) error {
	date, ok := value.(time.Time)
	if !ok {
		return errors.New("bad time.Time type assertion")
	}
	d.Time = date
	return nil
}

func (d Date) Value() (driver.Value, error) {
	return d.Time, nil
}

func (d *Date) MarshalJSON() ([]byte, error) {
	date := d.Time.Format("01/02")
	dayOfWeeks := [...]string{"日", "月", "火", "水", "木", "金", "土"}
	dayOfWeek := dayOfWeeks[d.Time.Weekday()]
	return []byte(`"` + date + `(` + dayOfWeek + `)` + `"`), nil
}

func (b BitBool) Value() (driver.Value, error) {
	if b {
		return []byte{1}, nil
	}
	return []byte{0}, nil
}

func (b *BitBool) Scan(src interface{}) error {
	bitBool, ok := src.([]byte)
	if !ok {
		return errors.New("bad []byte type assertion")
	}
	*b = bitBool[0] == 1
	return nil
}
