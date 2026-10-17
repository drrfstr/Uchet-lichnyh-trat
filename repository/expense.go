package repository

import (
	"database/sql"
	"errors"
	"log"
	"time"

	"project/models"

	_ "modernc.org/sqlite"
	"os"
	"path/filepath"
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

	return &ExpenseRepository{db:db}
}

func (r *ExpenseRepository) Create(amount int, description string, date time.Time) (int64, error) {
	result, err := r.db.Exec(
		"INSERT INTO expenses (amount, description, date) VALUES (?, ?, ?)",
		amount, description, date,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *ExpenseRepository) GetAll() ([]models.Expense, error) {
	rows, err := r.db.Query("SELECT id, amount, description, date FROM expenses ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var expenses []models.Expense
	for rows.Next() {
		var e models.Expense
		err := rows.Scan(&e.ID, &e.Amount, &e.Description, &e.Date)
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

func (r *ExpenseRepository) GetByID(id int64) (models.Expense, error) {
	var e models.Expense
	err := r.db.QueryRow(
		"SELECT id, amount, description, date FROM expenses WHERE id = ?",
		id,
	).Scan(&e.ID, &e.Amount, &e.Description, &e.Date)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Expense{}, errors.New("трата не найдена")
		}
		return models.Expense{}, err
	}
	return e, nil
}

func (r *ExpenseRepository) Update(id int64, amount int, description string, date time.Time) error {
	result, err := r.db.Exec(
		"UPDATE expenses SET amount = ?, description = ?, date = ? WHERE id = ?",
		amount, description, date, id,
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

func (r *ExpenseRepository) Delete(id int64) error {
	result, err := r.db.Exec("DELETE FROM expenses WHERE id = ?", id)
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

func (r *ExpenseRepository) Close() {
	r.db.Close()
}