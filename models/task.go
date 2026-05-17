package models

import (
	"time"

	"gorm.io/gorm"
)

type Status string

const (
	StatusPending    Status = "pending"
	StatusInProgress Status = "in_progress"
	StatusCompleted  Status = "completed"
)

type Task struct {
	gorm.Model
	UserID      uint       `json:"user_id" gorm:"not null"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      Status     `json:"status"`
	DueDate     *time.Time `json:"due_date"`
}

type CreateTaskInput struct {
	Title       string     `json:"title" binding:"required"`
	Description string     `json:"description"`
	DueDate     *time.Time `json:"due_date"`
}

type UpdateTaskInput struct {
	Title       *string    `json:"title"`
	Description *string    `json:"description"`
	Status      *Status    `json:"status"`
	DueDate     *time.Time `json:"due_date"`
}

type PatchTaskRequest struct {
	Title       *string    `json:"title"`
	Description *string    `json:"description"`
	Status      *Status    `json:"status"`
	DueDate     *time.Time `json:"due_date"`
}
