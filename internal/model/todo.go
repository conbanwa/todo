package model

import "time"

type Status string

const (
	NotStarted Status = "not_started"
	InProgress Status = "in_progress"
	Completed  Status = "completed"
)

type Todo struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	DueDate     time.Time `json:"due_date,omitempty"`
	Status      Status    `json:"status"`
	Priority    int       `json:"priority,omitempty"`
	Tags        []string  `json:"tags,omitempty"`
	TeamID      int64     `json:"team_id"`
}
