package repositories

import (
	"backend/models"
	"database/sql"
	"errors"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetAll() ([]models.User, error) {
	rows, err := r.db.Query(`
		SELECT id, first_name, last_name, email, phone, avatar_url, bio, points_balance, created_at, updated_at 
		FROM users
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		err := rows.Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.Phone, &u.AvatarURL, &u.Bio, &u.PointsBalance, &u.CreatedAt, &u.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

func (r *UserRepository) GetByID(id int64) (*models.User, error) {
	var u models.User
	err := r.db.QueryRow(`
		SELECT id, first_name, last_name, email, phone, avatar_url, bio, points_balance, created_at, updated_at 
		FROM users WHERE id = ?
	`, id).Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.Phone, &u.AvatarURL, &u.Bio, &u.PointsBalance, &u.CreatedAt, &u.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *UserRepository) Create(user *models.User) error {
	result, err := r.db.Exec(`
		INSERT INTO users (first_name, last_name, email, phone, avatar_url, bio, points_balance, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, user.FirstName, user.LastName, user.Email, user.Phone, user.AvatarURL, user.Bio, user.PointsBalance, user.CreatedAt, user.UpdatedAt)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	user.ID = id
	return nil
}

func (r *UserRepository) Update(id int64, user *models.User) error {
	result, err := r.db.Exec(`
		UPDATE users 
		SET first_name = ?, last_name = ?, email = ?, phone = ?, avatar_url = ?, bio = ?, updated_at = ?
		WHERE id = ?
	`, user.FirstName, user.LastName, user.Email, user.Phone, user.AvatarURL, user.Bio, user.UpdatedAt, id)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("user not found")
	}

	return nil
}

func (r *UserRepository) Delete(id int64) error {
	result, err := r.db.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("user not found")
	}

	return nil
}

func (r *UserRepository) UpdateBalance(tx *sql.Tx, userID int64, amount int64) error {
	_, err := tx.Exec("UPDATE users SET points_balance = points_balance + ?, updated_at = ? WHERE id = ?", amount, models.Now(), userID)
	return err
}

func (r *UserRepository) GetBalance(tx *sql.Tx, userID int64) (int64, error) {
	var balance int64
	err := tx.QueryRow("SELECT points_balance FROM users WHERE id = ?", userID).Scan(&balance)
	return balance, err
}

func (r *UserRepository) GetLastTransferRecipient(fromUserID int64) (int64, error) {
	var toUserID int64
	err := r.db.QueryRow(`
		SELECT to_user_id FROM transfers 
		WHERE from_user_id = ? AND status = 'completed'
		ORDER BY transfer_id DESC 
		LIMIT 1
	`, fromUserID).Scan(&toUserID)

	if err == sql.ErrNoRows {
		return 0, nil
	}

	return toUserID, err
}
