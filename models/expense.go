package models

import(
	"time"
)

type Expense struct {
	Amount 			int
	Description 	string
	Date 			time.Time
}

var Expense_list []Expense

func AddExpense(expense Expense) {
		Expense_list = append(Expense_list, expense)
}

func GetExpenses() []Expense {
		return Expense_list
}