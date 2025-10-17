# Code Style Guide

This project follows the style guidelines from [KBTG AI Workshop](https://github.com/mikelopster/kbtg-ai-workshop-oct).

## File Organization

### Clean Architecture Structure
```
backend/
├── models/          # Data structures and DTOs
├── repositories/    # Database access layer
├── services/        # Business logic layer
├── handlers/        # HTTP handlers (controllers)
├── main.go          # Application entry point
├── database.go      # Database initialization
└── error.go         # Error handling
```

## Naming Conventions

### Go Files
- Use descriptive names: `repository_user.go`, `service_transfer.go`
- Group by feature/domain
- Prefix by layer: `repository_`, `service_`, `handler_`

### Functions
- Use PascalCase for exported functions: `NewUserService()`, `CreateTransfer()`
- Use camelCase for private functions: `validateAmount()`, `checkBalance()`
- Clear action verbs: `Create`, `Update`, `Delete`, `Get`, `List`

### Variables
- Use camelCase: `userId`, `transferId`, `idemKey`
- Descriptive names: `fromBalance`, `toBalance`, `completedAt`
- Avoid abbreviations except common ones: `id`, `db`, `tx`

### Constants
- Use UPPER_SNAKE_CASE or PascalCase based on context
- Group related constants

## Package Structure

```go
// models/models.go
package models

import "time"

type User struct {
    ID            int64     `json:"id"`
    FirstName     string    `json:"first_name"`
    LastName      string    `json:"last_name"`
    PointsBalance int64     `json:"points_balance"`
    CreatedAt     time.Time `json:"created_at"`
    UpdatedAt     time.Time `json:"updated_at"`
}
```

## Repository Pattern

```go
// repositories/repository_user.go
package repositories

import (
    "backend/models"
    "database/sql"
)

type UserRepository struct {
    DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
    return &UserRepository{DB: db}
}

func (r *UserRepository) GetByID(id int64) (*models.User, error) {
    // Implementation
}
```

## Service Pattern

```go
// services/service_user.go
package services

import (
    "backend/models"
    "backend/repositories"
)

type UserService struct {
    repo *repositories.UserRepository
}

func NewUserService(repo *repositories.UserRepository) *UserService {
    return &UserService{repo: repo}
}

func (s *UserService) Create(req *models.CreateUserRequest) (*models.User, error) {
    // Validation
    // Business logic
    // Call repository
}
```

## Handler Pattern

```go
// handlers/handler_user.go
package handlers

import (
    "backend/models"
    "backend/services"
    "github.com/gofiber/fiber/v2"
)

type UserHandler struct {
    service *services.UserService
}

func NewUserHandler(service *services.UserService) *UserHandler {
    return &UserHandler{service: service}
}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
    var req models.CreateUserRequest
    if err := c.BodyParser(&req); err != nil {
        return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
    }
    
    user, err := h.service.Create(&req)
    if err != nil {
        return err
    }
    
    return c.Status(fiber.StatusCreated).JSON(user)
}
```

## Error Handling

```go
// Use fiber.NewError for HTTP errors
return fiber.NewError(fiber.StatusBadRequest, "invalid request")
return fiber.NewError(fiber.StatusNotFound, "user not found")
return fiber.NewError(fiber.StatusConflict, "insufficient balance")

// Use standard errors for internal logic
return errors.New("validation failed")
return fmt.Errorf("failed to create user: %w", err)
```

## JSON Tags

```go
// Use camelCase for API responses (matching OpenAPI spec)
type Transfer struct {
    TransferID  int64  `json:"transferId" db:"transfer_id"`
    FromUserID  int64  `json:"fromUserId" db:"from_user_id"`
    ToUserID    int64  `json:"toUserId" db:"to_user_id"`
    Amount      int64  `json:"amount" db:"amount"`
}

// Use snake_case for database columns
// Use db tags for SQL mapping
```

## Database Operations

```go
// Use transactions for multi-step operations
tx, err := db.Begin()
if err != nil {
    return err
}
defer tx.Rollback()

// ... perform operations ...

if err := tx.Commit(); err != nil {
    return err
}

// Use prepared statements and parameterized queries
result, err := tx.Exec(`
    INSERT INTO users (first_name, last_name, created_at, updated_at)
    VALUES (?, ?, ?, ?)
`, user.FirstName, user.LastName, now, now)
```

## Testing

```go
// Test function naming: Test + Feature + Condition
func TestUserNameValidation(t *testing.T) {}
func TestTransferAmountValidation(t *testing.T) {}
func TestNoConsecutiveSameRecipient(t *testing.T) {}

// Use table-driven tests
tests := []struct {
    name      string
    input     string
    wantError bool
}{
    {"Valid input", "test", false},
    {"Invalid input", "", true},
}

for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        // Test implementation
    })
}
```

## Comments

```go
// Package-level comments
// Package models provides data structures and DTOs for the transfer API

// Function comments for exported functions
// NewUserService creates a new user service with the given repository
func NewUserService(repo *repositories.UserRepository) *UserService {
    return &UserService{repo: repo}
}

// Inline comments for complex logic
// Check if user has sufficient balance before transfer
if fromBalance < req.Amount {
    return nil, errors.New("insufficient balance")
}
```

## Best Practices

### ✅ Do
- Use dependency injection
- Separate concerns (layers)
- Validate input at handler level
- Use transactions for multi-step operations
- Return errors, don't panic
- Use context for cancellation
- Close resources properly (defer)

### ❌ Don't
- Mix business logic in handlers
- Use global variables
- Ignore errors
- Use magic numbers (use constants)
- Commit database connection strings
- Log sensitive information

## API Design

### Follow OpenAPI Specification
- Use camelCase for JSON fields
- Include proper HTTP status codes
- Use appropriate HTTP methods (GET, POST, PUT, DELETE)
- Version your APIs if needed
- Document with OpenAPI/Swagger

### Response Format
```json
{
  "transfer": {
    "transferId": 1,
    "fromUserId": 1,
    "toUserId": 2,
    "amount": 100,
    "status": "completed"
  }
}
```

### Error Response Format
```json
{
  "error": "insufficient balance"
}
```

## References

- [KBTG AI Workshop Repository](https://github.com/mikelopster/kbtg-ai-workshop-oct)
- [Workshop 4 - Transfer API](https://github.com/mikelopster/kbtg-ai-workshop-oct/tree/main/workshop-4)
- [OpenAPI Specification](https://github.com/mikelopster/kbtg-ai-workshop-oct/blob/main/workshop-4/specs/transfer.yml)

---

**Last Updated:** October 17, 2025  
**Based on:** KBTG AI Workshop - Workshop 4 Guidelines
