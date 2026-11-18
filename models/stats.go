package models 

type CategoryStat struct {
	Name  string `json:"name"`
	Total int    `json:"total"`
}

type MonthStat struct {
	Month string `json:"month"`
	Total int    `json:"total"`
}