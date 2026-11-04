package repository

import (
	"database/sql"
	"errors"
	"project/models"
)

func (r *ExpenseRepository) CreateUser(u models.User) error {
	_, err := r.db.Exec(
		"INSERT INTO users (name, email, password_hash) VALUES (?, ?, ?)",
		u.Name, u.Email, u.PasswordHash,
	)
	return err
}

func (r *ExpenseRepository) GetUserByEmail(email string) (models.User, error) {
	var u models.User
	err := r.db.QueryRow(
		"SELECT id, name, email, password_hash FROM users WHERE email = ?", email,
	).Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, errors.New("пользователь не найден")
		}
		return models.User{}, err
	}
	return u, nil
}

func (r *ExpenseRepository) GetUserByID(id int) (models.User, error) {
 	var u models.User
 	err := r.db.QueryRow(
  		"SELECT id, name, email FROM users WHERE id = ?", id,
	).Scan(&u.ID, &u.Name, &u.Email)

 	if err != nil {
  		if errors.Is(err, sql.ErrNoRows) {
   			return models.User{}, errors.New("пользователь не найден")
 		}
  		return models.User{}, err
 	}
	return u, nil
}