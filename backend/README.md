# Points Transfer API

Go + Fiber backend with SQLite for managing users and point transfers.

## Features

- ✅ User CRUD operations
- ✅ Point transfer between users
- ✅ Idempotency support
- ✅ Transaction logging (point_ledger)
- ✅ Business rule validations:
  - User names limited to 3 characters
  - Transfer amount max 2.00 with 2 decimal places
  - No consecutive transfers to same recipient

## Tech Stack

- **Framework**: Go Fiber v2
- **Database**: SQLite3
- **Testing**: Go standard testing package

## Setup

### Prerequisites

- Go 1.21+
- SQLite3

### Installation

```bash
cd backend
go mod tidy
```

### Run

```bash
go run .
```

Server starts on `http://localhost:3000`

### Run Tests

```bash
go test -v
```

## API Endpoints

### Users

- `GET /api/users` - List all users
- `GET /api/users/:id` - Get user by ID
- `POST /api/users` - Create user
- `PUT /api/users/:id` - Update user
- `DELETE /api/users/:id` - Delete user

### Transfers

- `POST /api/transfers` - Create transfer
- `GET /api/transfers/:id` - Get transfer by ID

## API Examples

### Create User

```bash
curl -X POST http://localhost:3000/api/users \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "Tom",
    "last_name": "Lee",
    "email": "tom@example.com"
  }'
```

### Create Transfer

```bash
curl -X POST http://localhost:3000/api/transfers \
  -H "Content-Type: application/json" \
  -d '{
    "from_user_id": 1,
    "to_user_id": 2,
    "amount": 1.50,
    "idempotency_key": "unique-key-123",
    "note": "Payment for lunch"
  }'
```

## Business Rules

### User Validation
- `first_name` and `last_name` must not exceed 3 characters (enforced in service layer)
- Both fields are required

### Transfer Validation
1. **Amount Limits**: Maximum 2.00 per transfer, at most 2 decimal places
2. **No Consecutive Same Recipient**: Cannot transfer to the same user as the last completed transfer
3. **Idempotency**: Duplicate `idempotency_key` returns existing transfer
4. **Balance Check**: Sender must have sufficient balance
5. **User Validation**: Both sender and receiver must exist

## Database Schema

See [database.md](./database.md) for complete ER diagram.

### Key Tables

- **users**: User profiles and point balances
- **transfers**: Transfer records with idempotency
- **point_ledger**: Append-only transaction log

## Testing

Three main test suites:

1. **TestUserNameValidation**: Validates 3-character limit on names
2. **TestTransferAmountValidation**: Validates amount limits and decimals
3. **TestNoConsecutiveSameRecipient**: Validates no repeat recipients

Run with:
```bash
go test -v -run TestUserNameValidation
go test -v -run TestTransferAmountValidation
go test -v -run TestNoConsecutiveSameRecipient
```

## Architecture

```
backend/
├── main.go                  # Entry point
├── database.go              # DB initialization & migrations
├── models.go                # Data structures
├── repository_user.go       # User data access
├── repository_transfer.go   # Transfer data access
├── repository_ledger.go     # Ledger data access
├── service_user.go          # User business logic
├── service_transfer.go      # Transfer business logic
├── handler_user.go          # User HTTP handlers
├── handler_transfer.go      # Transfer HTTP handlers
├── errors.go                # Error handling
└── main_test.go             # Unit tests
```

## API Documentation

- OpenAPI spec: [result.yml](./result.yml)
- Sequence diagrams: [result.md](./result.md)
- Database schema: [database.md](./database.md)

## License

MIT
