package queue

import "time"

type Status string

const (
	StatusPending   Status = "pending"
	StatusRunning   Status = "running"
	StatusDone      Status = "done"
	StatusFailed    Status = "failed"
	StatusCancelled Status = "cancelled"
)

// QueueRecord is the DB row in the `queues` table
type QueueRecord struct {
	ID        int64     `gorm:"primayKey;autoIncrement" json:"id"`
	UniqueID  string    `gorm:"column:unique_id" json:"unique_id"`
	Queue     string    `gorm:"column:queue" json:"queue"`
	Payload   string    `gorm:"column:payload" json:"payload"`
	Status    Status    `gorm:"column:status" json:"status"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime json:"updated_at"`
}

func (queueRecord *QueueRecord) TableName() string {
	return "queues"
}

type QueueMessage struct {
	UniqueID string                 `json:"unique_id"`
	Payload  map[string]interface{} `json:"payload"`
}
