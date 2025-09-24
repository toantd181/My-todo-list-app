package models

import (
	"time"
)

type Todo struct {
	ID         int       `json:"id" db:"id"`
	Name       string    `json:"name" db:"name" binding:"required, max = 500"`
	CreateTime string    `json: "create_time" db:"create_time"`
	Importance int       `json:"importance" db:"importance"`
	Flag       bool      `json: "flag" db:"flag"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json: "updated_at" db: "updated_at"`
}

type CreateTodoRequest struct {
	Name       string `json: "name" binding: "required, max = 500"`
	CreateTime string `json: "create_time, omitempty"`
	Importance int    `json: "importance, omitempty"`
}

type UpdateTodoRequest struct {
	Name       *string `json: "name, omitempty"`
	CreateTime *string `json: "create_time, omitempty"`
	Importance *int    `json: "importance, omitempty"`
	Flag       *bool   `json: "flag, omitempty"`
}

type ReorderRequest struct {
	TodoIDs []int `json: "todo_ids" binding: "required"`
}
