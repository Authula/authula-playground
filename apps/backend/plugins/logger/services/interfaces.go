package services

import (
	"context"

	"github.com/Authula/authula-playground/plugins/logger/types"
)

type LoggerService interface {
	CreateLogEntry(ctx context.Context, eventType string, details map[string]any) (*types.LogEntry, error)
	GetLogEntry(ctx context.Context, id int64) (*types.LogEntry, error)
	GetAllLogs(ctx context.Context) ([]types.LogEntry, error)
	DeleteLogEntry(ctx context.Context, id int64) error
	GetLogCount(ctx context.Context) (int64, error)
	HasReachedMaxLogs(ctx context.Context) (bool, error)
}
