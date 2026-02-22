package database

import (
	"context"
	"encoding/json"
	"time"

	"github.com/charmbracelet/log"
	"github.com/kiss2u/SaveAny-Bot/database/sqlite"
	"gorm.io/gorm"
)

// TaskState represents the persisted state of a task
type TaskState struct {
	ID          string    `gorm:"primaryKey" json:"id"`
	Title       string    `json:"title"`
	Type        string    `json:"type"`
	Status      string    `json:"status"` // pending, running, completed, failed, cancelled
	Data        string    `json:"data"`   // JSON serialized task data
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	Error       string    `json:"error,omitempty"`
}

func (TaskState) TableName() string {
	return "task_states"
}

func InitTaskState(ctx context.Context) error {
	return GetDB(ctx).AutoMigrate(&TaskState{})
}

func SaveTaskState(ctx context.Context, task *TaskState) error {
	return GetDB(ctx).Save(task).Error
}

func GetTaskState(ctx context.Context, id string) (*TaskState, error) {
	var task TaskState
	err := GetDB(ctx).Where("id = ?", id).First(&task).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func GetPendingTasks(ctx context.Context) ([]TaskState, error) {
	var tasks []TaskState
	err := GetDB(ctx).Where("status IN ?", []string{"pending", "running"}).Find(&tasks).Error
	return tasks, err
}

func DeleteTaskState(ctx context.Context, id string) error {
	return GetDB(ctx).Where("id = ?", id).Delete(&TaskState{}).Error
}

// PersistTask saves a task to the database for recovery
func PersistTask(ctx context.Context, id, title, taskType, data string) error {
	task := &TaskState{
		ID:        id,
		Title:     title,
		Type:      taskType,
		Status:    "pending",
		Data:      data,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return SaveTaskState(ctx, task)
}

// UpdateTaskStatus updates the status of a persisted task
func UpdateTaskStatus(ctx context.Context, id, status, errMsg string) error {
	updates := map[string]interface{}{
		"status":     status,
		"updated_at": time.Now(),
	}
	if errMsg != "" {
		updates["error"] = errMsg
	}
	if status == "completed" || status == "failed" || status == "cancelled" {
		now := time.Now()
		updates["completed_at"] = &now
	}
	return GetDB(ctx).Model(&TaskState{}).Where("id = ?", id).Updates(updates).Error
}
