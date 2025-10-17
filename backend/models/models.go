package models

import "time"

type User struct {
	ID            int64     `json:"id"`
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	Email         string    `json:"email,omitempty"`
	Phone         string    `json:"phone,omitempty"`
	AvatarURL     string    `json:"avatar_url,omitempty"`
	Bio           string    `json:"bio,omitempty"`
	PointsBalance int64     `json:"points_balance"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Transfer struct {
	TransferID  int64      `json:"transferId" db:"transfer_id"`
	IdemKey     string     `json:"idemKey" db:"idempotency_key"`
	FromUserID  int64      `json:"fromUserId" db:"from_user_id"`
	ToUserID    int64      `json:"toUserId" db:"to_user_id"`
	Amount      int64      `json:"amount" db:"amount"`
	Status      string     `json:"status" db:"status"`
	Note        string     `json:"note" db:"note"`
	CreatedAt   time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time  `json:"updatedAt" db:"updated_at"`
	CompletedAt *time.Time `json:"completedAt,omitempty" db:"completed_at"`
	FailReason  *string    `json:"failReason,omitempty" db:"fail_reason"`
}

type PointLedger struct {
	ID           int64     `json:"id"`
	UserID       int64     `json:"user_id"`
	Change       int64     `json:"change"`
	BalanceAfter int64     `json:"balance_after"`
	EventType    string    `json:"event_type"`
	TransferID   *int64    `json:"transfer_id,omitempty"`
	Reference    string    `json:"reference,omitempty"`
	Metadata     string    `json:"metadata,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

type CreateUserRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email,omitempty"`
	Phone     string `json:"phone,omitempty"`
	AvatarURL string `json:"avatar_url,omitempty"`
	Bio       string `json:"bio,omitempty"`
}

type UpdateUserRequest struct {
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Email     string `json:"email,omitempty"`
	Phone     string `json:"phone,omitempty"`
	AvatarURL string `json:"avatar_url,omitempty"`
	Bio       string `json:"bio,omitempty"`
}

type CreateTransferRequest struct {
	FromUserID int64  `json:"fromUserId"`
	ToUserID   int64  `json:"toUserId"`
	Amount     int64  `json:"amount"` // in cents/points
	Note       string `json:"note,omitempty"`
}

type TransferListQuery struct {
	UserID   int64 `query:"userId"`
	Page     int   `query:"page"`
	PageSize int   `query:"pageSize"`
}

type TransferListResponse struct {
	Data     []Transfer `json:"data"`
	Page     int        `json:"page"`
	PageSize int        `json:"pageSize"`
	Total    int        `json:"total"`
}

func Now() time.Time {
	return time.Now().UTC()
}
