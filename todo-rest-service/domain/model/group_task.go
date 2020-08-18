package model

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

type GroupTasksListForEachUser struct {
	GroupTasksUsersList []GroupTasksUser `json:"group_tasks_list_for_each_user"`
}

type GroupTasksUser struct {
	ID        int         `json:"id"          db:"id"`
	UserID    string      `json:"user_id"     db:"user_id"`
	GroupID   int         `json:"group_id"    db:"group_id"`
	TasksList []GroupTask `json:"tasks_list"`
}

type GroupTasksList struct {
	GroupTasksList []GroupTask `json:"group_tasks_list"`
}

type GroupTask struct {
	ID               int        `json:"id"                   db:"id"`
	BaseDate         NullTime   `json:"base_date"            db:"base_date"`
	CycleType        NullString `json:"cycle_type"           db:"cycle_type"`
	Cycle            NullInt    `json:"cycle"                db:"cycle"`
	TaskName         string     `json:"task_name"            db:"task_name"`
	GroupID          int        `json:"group_id"             db:"group_id"`
	GroupTasksUserID NullInt    `json:"group_tasks_users_id" db:"group_tasks_users_id"`
}

func NewGroupTasksListForEachUser(groupTasksUsersList []GroupTasksUser) GroupTasksListForEachUser {
	return GroupTasksListForEachUser{
		GroupTasksUsersList: groupTasksUsersList,
	}
}

func NewGroupTasksList(groupTasksList []GroupTask) GroupTasksList {
	return GroupTasksList{
		GroupTasksList: groupTasksList,
	}
}

type NullTime struct {
	sql.NullTime
}

type NullString struct {
	sql.NullString
}

type NullInt struct {
	Int   int
	Valid bool
}

func (nt *NullTime) MarshalJSON() ([]byte, error) {
	if !nt.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(nt.Time)
}

func (nt *NullTime) UnmarshalJSON(b []byte) error {
	if bytes.Equal(b, []byte("null")) {
		return nil
	}
	trimData := strings.Trim(string(b), "\"")[:10]
	dateTime, err := time.Parse("2006-01-02", trimData)
	if err != nil {
		return err
	}
	nt.Time, nt.Valid = dateTime, true
	return nil
}

func (ns *NullString) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ns.String)
}

func (ns *NullString) UnmarshalJSON(b []byte) error {
	if bytes.Equal(b, []byte("null")) {
		return nil
	}
	if err := json.Unmarshal(b, &ns.String); err != nil {
		return err
	}
	ns.Valid = true
	return nil
}

func (ni *NullInt) MarshalJSON() ([]byte, error) {
	if !ni.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ni.Int)
}

func (ni *NullInt) UnmarshalJSON(b []byte) error {
	if bytes.Equal(b, []byte("null")) {
		return nil
	}
	if err := json.Unmarshal(b, &ni.Int); err != nil {
		return err
	}
	ni.Valid = true
	return nil
}

func (ni *NullInt) Scan(value interface{}) error {
	if value == nil {
		ni.Int, ni.Valid = 0, false
		return nil
	}

	intValue, ok := value.(int64)
	if !ok {
		return errors.New("type assertion error")
	}
	ni.Int, ni.Valid = int(intValue), true
	return nil
}

func (ni NullInt) Value() (driver.Value, error) {
	if !ni.Valid {
		return nil, nil
	}
	return int64(ni.Int), nil
}
