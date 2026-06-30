package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	ID          int64      `json:"id"`
	ServiceName string     `json:"service_name"`
	Price       int        `json:"price"`
	UserID      uuid.UUID  `json:"user_id"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     *time.Time `json:"end_date,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type SubscriptionRepository struct {
	db *sql.DB
}

func NewSubscriptionRepository(db *sql.DB) *SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

type CreateSubscriptionRequest struct {
	UserID      string `json:"user_id"`
	ServiceName string `json:"service_name"`
	Price       int    `json:"price"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
}

func (r *SubscriptionRepository) GetByID(ctx context.Context, id int64) (*Subscription, error) {
	const query = `
		SELECT 
			id,
			service_name,
			price,
			user_id,
			start_date,
			end_date,
			created_at,
			updated_at
		FROM subscriptions
		WHERE id = $1
	`

	var sub Subscription

	var endDateNull sql.NullTime

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&sub.ID,
		&sub.ServiceName,
		&sub.Price,
		&sub.UserID,
		&sub.StartDate,
		&endDateNull,
		&sub.CreatedAt,
		&sub.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("not found")
	}
	if err != nil {
		return nil, err
	}

	// Конвертируем sql.NullTime в *time.Time
	if endDateNull.Valid {
		sub.EndDate = &endDateNull.Time
	} else {
		sub.EndDate = nil
	}

	return &sub, nil
}


func (r *SubscriptionRepository) Create(ctx context.Context, userIDStr, serviceName string, price int, startDateStr, endDateStr string) (*Subscription, error) {
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid user_id: %w", err)
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		return nil, fmt.Errorf("invalid start_date: %w", err)
	}

	var endDate *time.Time
	if endDateStr != "" {
		t, err := time.Parse("2006-01-02", endDateStr)
		if err != nil {
			return nil, fmt.Errorf("invalid end_date: %w", err)
		}
		endDate = &t
	}

	query := `
		INSERT INTO subscriptions (user_id, service_name, price, start_date, end_date)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`

	var sub Subscription
	err = r.db.QueryRowContext(ctx, query, userID, serviceName, price, startDate, endDate).Scan(
		&sub.ID, &sub.CreatedAt, &sub.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	sub.UserID = userID
	sub.ServiceName = serviceName
	sub.Price = price
	sub.StartDate = startDate
	sub.EndDate = endDate

	return &sub, nil
}
