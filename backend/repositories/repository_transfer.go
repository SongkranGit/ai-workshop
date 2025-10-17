package repositories

import (
	"backend/models"
	"database/sql"
	"errors"
)

type TransferRepository struct {
	DB *sql.DB
}

func NewTransferRepository(db *sql.DB) *TransferRepository {
	return &TransferRepository{DB: db}
}

func (r *TransferRepository) GetByIdemKey(key string) (*models.Transfer, error) {
	var t models.Transfer
	err := r.DB.QueryRow(`
		SELECT transfer_id, idempotency_key, from_user_id, to_user_id, amount, status, note, created_at, updated_at, completed_at, fail_reason
		FROM transfers WHERE idempotency_key = ?
	`, key).Scan(&t.TransferID, &t.IdemKey, &t.FromUserID, &t.ToUserID, &t.Amount, &t.Status, &t.Note, &t.CreatedAt, &t.UpdatedAt, &t.CompletedAt, &t.FailReason)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (r *TransferRepository) GetByID(id int64) (*models.Transfer, error) {
	var t models.Transfer
	err := r.DB.QueryRow(`
		SELECT transfer_id, idempotency_key, from_user_id, to_user_id, amount, status, note, created_at, updated_at, completed_at, fail_reason
		FROM transfers WHERE transfer_id = ?
	`, id).Scan(&t.TransferID, &t.IdemKey, &t.FromUserID, &t.ToUserID, &t.Amount, &t.Status, &t.Note, &t.CreatedAt, &t.UpdatedAt, &t.CompletedAt, &t.FailReason)

	if err == sql.ErrNoRows {
		return nil, errors.New("transfer not found")
	}
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (r *TransferRepository) GetByUserID(userID int64, page, pageSize int) ([]models.Transfer, int, error) {
	offset := (page - 1) * pageSize

	// Get total count
	var total int
	err := r.DB.QueryRow(`
		SELECT COUNT(*) FROM transfers 
		WHERE from_user_id = ? OR to_user_id = ?
	`, userID, userID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get paginated data
	rows, err := r.DB.Query(`
		SELECT transfer_id, idempotency_key, from_user_id, to_user_id, amount, status, note, created_at, updated_at, completed_at, fail_reason
		FROM transfers 
		WHERE from_user_id = ? OR to_user_id = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`, userID, userID, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var transfers []models.Transfer
	for rows.Next() {
		var t models.Transfer
		err := rows.Scan(&t.TransferID, &t.IdemKey, &t.FromUserID, &t.ToUserID, &t.Amount, &t.Status, &t.Note, &t.CreatedAt, &t.UpdatedAt, &t.CompletedAt, &t.FailReason)
		if err != nil {
			return nil, 0, err
		}
		transfers = append(transfers, t)
	}

	return transfers, total, nil
}

func (r *TransferRepository) Create(tx *sql.Tx, transfer *models.Transfer) error {
	result, err := tx.Exec(`
		INSERT INTO transfers (idempotency_key, from_user_id, to_user_id, amount, status, note, created_at, updated_at, completed_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, transfer.IdemKey, transfer.FromUserID, transfer.ToUserID, transfer.Amount, transfer.Status, transfer.Note, transfer.CreatedAt, transfer.UpdatedAt, transfer.CompletedAt)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	transfer.TransferID = id
	return nil
}
