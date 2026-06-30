package models

import (
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	ID          int64
	UserID      uuid.UUID
	ServiceName string
	Price       int
	StartDate   time.Time
	EndDate     *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
