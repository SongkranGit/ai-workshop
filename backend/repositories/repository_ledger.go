package repositories

import (
	"backend/models"
	"database/sql"
)

type LedgerRepository struct {
	db *sql.DB
}

func NewLedgerRepository(db *sql.DB) *LedgerRepository {
	return &LedgerRepository{db: db}
}

func (r *LedgerRepository) Create(tx *sql.Tx, ledger *models.PointLedger) error {
	result, err := tx.Exec(`
		INSERT INTO point_ledger (user_id, change, balance_after, event_type, transfer_id, reference, metadata, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, ledger.UserID, ledger.Change, ledger.BalanceAfter, ledger.EventType, ledger.TransferID, ledger.Reference, ledger.Metadata, ledger.CreatedAt)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	ledger.ID = id
	return nil
}
