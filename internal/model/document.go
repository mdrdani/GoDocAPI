package model

import (
	"time"

	"github.com/google/uuid"
)

type Document struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Filename    string    `json:"filename" db:"filename"`
	StoragePath string    `json:"storage_path" db:"storage_path"`
	Size        int64     `json:"size" db:"size"`
	ContentType string    `json:"content_type" db:"content_type"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}
