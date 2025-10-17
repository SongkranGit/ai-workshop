package services

import (
	"backend/models"
	"backend/repositories"
	"errors"

	"github.com/google/uuid"
)

type TransferService struct {
	transferRepo *repositories.TransferRepository
	ledgerRepo   *repositories.LedgerRepository
	userRepo     *repositories.UserRepository
}

func NewTransferService(transferRepo *repositories.TransferRepository, ledgerRepo *repositories.LedgerRepository, userRepo *repositories.UserRepository) *TransferService {
	return &TransferService{
		transferRepo: transferRepo,
		ledgerRepo:   ledgerRepo,
		userRepo:     userRepo,
	}
}

func (s *TransferService) CreateTransfer(req *models.CreateTransferRequest) (*models.Transfer, error) {
	// Generate idempotency key (idemKey)
	idemKey := uuid.New().String()

	// Validation: Amount must be > 0
	if req.Amount <= 0 {
		return nil, errors.New("amount must be greater than 0")
	}

	// Validation: Cannot transfer to same recipient as last transfer
	lastRecipient, err := s.userRepo.GetLastTransferRecipient(req.FromUserID)
	if err != nil {
		return nil, err
	}
	if lastRecipient != 0 && lastRecipient == req.ToUserID {
		return nil, errors.New("cannot transfer to the same recipient as your last transfer")
	}

	// Validate users exist
	fromUser, err := s.userRepo.GetByID(req.FromUserID)
	if err != nil {
		return nil, errors.New("from_user not found")
	}

	_, err = s.userRepo.GetByID(req.ToUserID)
	if err != nil {
		return nil, errors.New("to_user not found")
	}

	if req.FromUserID == req.ToUserID {
		return nil, errors.New("cannot transfer to yourself")
	}

	// Check balance
	if fromUser.PointsBalance < req.Amount {
		return nil, errors.New("insufficient balance")
	}

	// Begin transaction
	db := s.transferRepo.DB
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	now := models.Now()
	completedAt := now
	transfer := &models.Transfer{
		IdemKey:     idemKey,
		FromUserID:  req.FromUserID,
		ToUserID:    req.ToUserID,
		Amount:      req.Amount,
		Status:      "completed",
		Note:        req.Note,
		CreatedAt:   now,
		UpdatedAt:   now,
		CompletedAt: &completedAt,
	}

	// Create transfer record
	err = s.transferRepo.Create(tx, transfer)
	if err != nil {
		return nil, err
	}

	// Update balances
	err = s.userRepo.UpdateBalance(tx, req.FromUserID, -req.Amount)
	if err != nil {
		return nil, err
	}

	err = s.userRepo.UpdateBalance(tx, req.ToUserID, req.Amount)
	if err != nil {
		return nil, err
	}

	// Get updated balances
	fromBalance, err := s.userRepo.GetBalance(tx, req.FromUserID)
	if err != nil {
		return nil, err
	}

	toBalance, err := s.userRepo.GetBalance(tx, req.ToUserID)
	if err != nil {
		return nil, err
	}

	// Create ledger entries
	fromLedger := &models.PointLedger{
		UserID:       req.FromUserID,
		Change:       -req.Amount,
		BalanceAfter: fromBalance,
		EventType:    "transfer_out",
		TransferID:   &transfer.TransferID,
		CreatedAt:    now,
	}
	err = s.ledgerRepo.Create(tx, fromLedger)
	if err != nil {
		return nil, err
	}

	toLedger := &models.PointLedger{
		UserID:       req.ToUserID,
		Change:       req.Amount,
		BalanceAfter: toBalance,
		EventType:    "transfer_in",
		TransferID:   &transfer.TransferID,
		CreatedAt:    now,
	}
	err = s.ledgerRepo.Create(tx, toLedger)
	if err != nil {
		return nil, err
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return transfer, nil
}

func (s *TransferService) GetByIdemKey(idemKey string) (*models.Transfer, error) {
	transfer, err := s.transferRepo.GetByIdemKey(idemKey)
	if err != nil {
		return nil, err
	}
	if transfer == nil {
		return nil, errors.New("transfer not found")
	}
	return transfer, nil
}

func (s *TransferService) GetByUserID(userID int64, page, pageSize int) (*models.TransferListResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 200 {
		pageSize = 20
	}

	transfers, total, err := s.transferRepo.GetByUserID(userID, page, pageSize)
	if err != nil {
		return nil, err
	}

	if transfers == nil {
		transfers = []models.Transfer{}
	}

	return &models.TransferListResponse{
		Data:     transfers,
		Page:     page,
		PageSize: pageSize,
		Total:    total,
	}, nil
}
