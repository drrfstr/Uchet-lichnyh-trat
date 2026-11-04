package repository

import (
	"database/sql"
	"errors"
	"log"
	"strings"

	// "time"

	"project/models"

	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

type ExpenseRepository struct {
	db *sql.DB
}

func NewExpenseRepository(dbPath string) *ExpenseRepository {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatal("не удалось открыть базу:", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("база не отвечает:", err)
	}

	dirPath := "./repository/migrations"
	dirEntries, err := os.ReadDir(dirPath)
	if err != nil {
		log.Fatal("ошибка чтения папки базы данных: ", err)
	}

	for _, entry := range dirEntries {

		if entry.IsDir() {
			continue
		}

		filePath := filepath.Join(dirPath, entry.Name())
		content, err := os.ReadFile(filePath)
		if err != nil {
			log.Printf("не удалось прочитать файл %s: %v", filePath, err)
			continue
		}

		sql := string(content)

		_, err = db.Exec(sql)
		if err != nil {
			log.Fatal("не удалось создать таблицу:", err)
		}

	}

	return &ExpenseRepository{db: db}
}

func (r *ExpenseRepository) Close() {
	r.db.Close()
}




func (r *ExpenseRepository) CreateExpense(e models.Expense) (int64, error) {
	result, err := r.db.Exec(
		"INSERT INTO expenses (amount, description, date, category_id, user_id) VALUES (?, ?, ?, ?, ?)",
		e.Amount, e.Description, e.Date, e.CategoryID, e.UserID,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *ExpenseRepository) GetAll() ([]models.Expense, error) {
	rows, err := r.db.Query(`
	SELECT e.id, e.amount, e.description, e.date, e.category_id, c.name as category_name, e.user_id
	FROM expenses e 
	LEFT JOIN categories c ON e.category_id = c.id 
	ORDER BY e.id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var expenses []models.Expense
	for rows.Next() {
		var e models.Expense
		err := rows.Scan(&e.ID, &e.Amount, &e.Description, &e.Date, &e.CategoryID, &e.CategoryName, &e.UserID)
		if err != nil {
			return nil, err
		}
		expenses = append(expenses, e)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return expenses, nil
}

func (r *ExpenseRepository) GetByID(id int64, userID int) (models.Expense, error) {
	var e models.Expense
	err := r.db.QueryRow(
		`SELECT e.id, e.amount, e.description, e.date, e.category_id, c.name as category_name, e.user_id
	FROM expenses e 
	LEFT JOIN categories c ON e.category_id = c.id
	WHERE e.id = ? AND user_id = ?`,
		id, userID,
	).Scan(&e.ID, &e.Amount, &e.Description, &e.Date, &e.CategoryID, &e.CategoryName, &e.UserID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Expense{}, errors.New("трата не найдена")
		}
		return models.Expense{}, err
	}
	return e, nil
}

func (r *ExpenseRepository) Update(e models.Expense) error {
	result, err := r.db.Exec(
		"UPDATE expenses SET amount = ?, description = ?, date = ?, category_id = ? WHERE id = ? AND user_id = ?",
		e.Amount, e.Description, e.Date, e.CategoryID, e.ID, e.UserID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("трата не найдена")
	}
	return nil
}

func (r *ExpenseRepository) Delete(id int64, userID int) error {
	result, err := r.db.Exec("DELETE FROM expenses WHERE id = ? AND user_id = ?", id, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("трата не найдена")
	}
	return nil
}




func (r *ExpenseRepository) GetAllCategories() ([]models.Category, error) {
	rows, err := r.db.Query("SELECT id, name FROM categories ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var c models.Category
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, rows.Err()
}

func (r *ExpenseRepository) GetAllFiltered(categoryID int, dateFrom, dateTo string, userID int) ([]models.Expense, error) {
	query := `
		SELECT e.id, e.amount, e.description, e.date, e.category_id, c.name AS category_name, e.user_id
		FROM expenses e 
		LEFT JOIN categories c ON e.category_id = c.id
	`
	
	conditions := []string{"e.user_id = ?"}
	args := []any{userID}

	if categoryID > 0 {
		conditions = append(conditions, "e.category_id = ?")
		args = append(args, categoryID)
	}
	if dateFrom != "" {
		conditions = append(conditions, "e.date >= ?")
		args = append(args, dateFrom)
	}
	if dateTo != "" {
		conditions = append(conditions, "e.date <= ?")
		args = append(args, dateTo)
	}


	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY e.id"


	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var expenses []models.Expense
	for rows.Next() {
		var e models.Expense
		err := rows.Scan(&e.ID, &e.Amount, &e.Description, &e.Date, &e.CategoryID, &e.CategoryName, &e.UserID)
		if err != nil {
			return nil, err
		}
		expenses = append(expenses, e)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return expenses, nil
}