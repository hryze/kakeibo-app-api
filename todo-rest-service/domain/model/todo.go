package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

type TodoList struct {
	ImplementationTodoList []Todo `json:"implementation_todo_list"`
	DueTodoList            []Todo `json:"due_todo_list"`
}

type ExpiredTodoList struct {
	ExpiredTodoList []Todo `json:"expired_todo_list"`
}

type SearchTodoList struct {
	SearchTodoList []Todo `json:"search_todo_list"`
}

type Todo struct {
	ID                 int       `json:"id"                  db:"id"`
	PostedDate         time.Time `json:"posted_date"         db:"posted_date"`
	UpdatedDate        time.Time `json:"updated_date"        db:"updated_date"`
	ImplementationDate Date      `json:"implementation_date" db:"implementation_date" validate:"required,date"`
	DueDate            Date      `json:"due_date"            db:"due_date"            validate:"required,date"`
	TodoContent        string    `json:"todo_content"        db:"todo_content"        validate:"required,max=100,blank"`
	CompleteFlag       BitBool   `json:"complete_flag"       db:"complete_flag"`
}

type Date struct {
	time.Time
}

type BitBool bool

func NewTodoList(implementationTodoList []Todo, dueTodoList []Todo) TodoList {
	return TodoList{
		ImplementationTodoList: implementationTodoList,
		DueTodoList:            dueTodoList,
	}
}

func NewSearchTodoList(searchTodoList []Todo) SearchTodoList {
	return SearchTodoList{
		SearchTodoList: searchTodoList,
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
	date := d.Time.Format("2006/01/02")
	dayOfWeeks := [...]string{"日", "月", "火", "水", "木", "金", "土"}
	dayOfWeek := dayOfWeeks[d.Time.Weekday()]

	return []byte(`"` + date + `(` + dayOfWeek + `)` + `"`), nil
}

func (d *Date) UnmarshalJSON(data []byte) error {
	trimData := strings.Trim(string(data), "\"")[:10]
	format := trimData[4:5]
	var date time.Time
	var err error

	if format == "-" {
		date, err = time.Parse("2006-01-02", trimData)
		if err != nil {
			return err
		}
	} else if format == "/" {
		date, err = time.Parse("2006/01/02", trimData)
		if err != nil {
			return err
		}
	}

	d.Time = date

	return nil
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

func (t Todo) ShowTodo() (string, error) {
	b, err := json.Marshal(t)
	if err != nil {
		return string(b), err
	}

	return string(b), nil
}
