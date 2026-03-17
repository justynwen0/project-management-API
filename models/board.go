package models

import "time"

type Board struct {
	InternalID    int64     `json:"internal_id" db:"internal_id" gorm:"primaryKey"`
	PublicID      string    `json:"public_id" db:"public_id"`
	Title         string    `json:"title" db:"title"`
	Description   string    `json:"description" db:"description"`
	OwnerID       int64     `json:"owner_internal_id" db:"owner_internal_id" gorm:"column:inter_internal_id"`
	OwnerPublicID string    `json:"owner_public_id" db:"owner_public_id"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	DueDate       time.Time `json:"due_date,omitempty" db:"due_date"`
}
