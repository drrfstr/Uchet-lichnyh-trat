package models

import (
	"time"
)

type Expense struct {
	ID          int64     `json:"id"`
	Amount      int       `json:"amount"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
}