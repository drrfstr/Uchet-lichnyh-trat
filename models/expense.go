package models

import (
	"time"
)

type Expense struct {
	ID          int64     `json:"id"`
	Amount      int       `json:"amount"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
	CategoryID   int     `json:"category_id"`
	CategoryName string    `json:"category_name"`
}