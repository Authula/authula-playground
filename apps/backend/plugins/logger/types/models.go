package types

import (
	"time"

	"github.com/uptrace/bun"
)

type LogEntry struct {
	bun.BaseModel `bun:"table:log_entries"`

	ID        int64          `json:"id" bun:"column:id,pk,autoincrement"`
	EventType string         `json:"event_type" bun:"column:event_type"`
	Details   map[string]any `json:"details" bun:"column:details"`
	CreatedAt time.Time      `json:"created_at" bun:"column:created_at,default:current_timestamp"`
}
