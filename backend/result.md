# Transfer API Sequence Diagram

## POST /api/transfers - Create Transfer

```mermaid
sequenceDiagram
    actor Client
    participant API
    participant TransferService
    participant UserRepo
    participant TransferRepo
    participant LedgerRepo
    participant Database

    Client->>API: POST /api/transfers
    API->>TransferService: CreateTransfer(request)
    
    Note over TransferService: Validation 1: Check idempotency
    TransferService->>TransferRepo: GetByIdempotencyKey(key)
    TransferRepo->>Database: SELECT * FROM transfers WHERE idempotency_key=?
    Database-->>TransferRepo: result
    TransferRepo-->>TransferService: existing transfer or nil
    
    alt Transfer already exists
        TransferService-->>API: return existing transfer
        API-->>Client: 200 OK (idempotent)
    end
    
    Note over TransferService: Validation 2: Amount limits
    alt Amount > 2.00 or decimals > 2
        TransferService-->>API: error: invalid amount
        API-->>Client: 400 Bad Request
    end
    
    Note over TransferService: Validation 3: Check last recipient
    TransferService->>UserRepo: GetLastTransferRecipient(fromUserID)
    UserRepo->>Database: SELECT to_user_id FROM transfers<br/>WHERE from_user_id=? ORDER BY created_at DESC LIMIT 1
    Database-->>UserRepo: lastRecipientID
    UserRepo-->>TransferService: lastRecipientID
    
    alt Same recipient as last transfer
        TransferService-->>API: error: cannot transfer to same recipient
        API-->>Client: 400 Bad Request
    end
    
    Note over TransferService: Validate users exist
    TransferService->>UserRepo: GetByID(fromUserID)
    UserRepo->>Database: SELECT * FROM users WHERE id=?
    Database-->>UserRepo: fromUser
    UserRepo-->>TransferService: fromUser
    
    TransferService->>UserRepo: GetByID(toUserID)
    UserRepo->>Database: SELECT * FROM users WHERE id=?
    Database-->>UserRepo: toUser
    UserRepo-->>TransferService: toUser
    
    alt User not found
        TransferService-->>API: error: user not found
        API-->>Client: 404 Not Found
    end
    
    alt Insufficient balance
        TransferService-->>API: error: insufficient balance
        API-->>Client: 400 Bad Request
    end
    
    Note over TransferService,Database: Begin Transaction
    TransferService->>Database: BEGIN TRANSACTION
    
    TransferService->>TransferRepo: Create(tx, transfer)
    TransferRepo->>Database: INSERT INTO transfers
    Database-->>TransferRepo: transferID
    
    TransferService->>UserRepo: UpdateBalance(tx, fromUserID, -amount)
    UserRepo->>Database: UPDATE users SET points_balance -= amount
    
    TransferService->>UserRepo: UpdateBalance(tx, toUserID, +amount)
    UserRepo->>Database: UPDATE users SET points_balance += amount
    
    TransferService->>UserRepo: GetBalance(tx, fromUserID)
    UserRepo->>Database: SELECT points_balance FROM users
    Database-->>UserRepo: fromBalance
    
    TransferService->>UserRepo: GetBalance(tx, toUserID)
    UserRepo->>Database: SELECT points_balance FROM users
    Database-->>UserRepo: toBalance
    
    TransferService->>LedgerRepo: Create(tx, fromLedger)
    LedgerRepo->>Database: INSERT INTO point_ledger (transfer_out)
    
    TransferService->>LedgerRepo: Create(tx, toLedger)
    LedgerRepo->>Database: INSERT INTO point_ledger (transfer_in)
    
    TransferService->>Database: COMMIT TRANSACTION
    
    TransferService-->>API: transfer
    API-->>Client: 201 Created
```

## Key Validation Points

### 1. Name Length Validation (Users)
- Checked in `UserService.Create()` and `UserService.Update()`
- Uses `utf8.RuneCountInString()` to count characters (not bytes)
- Returns error if first_name or last_name > 3 characters

### 2. Transfer Amount Validation
- Maximum: 2.00 (200 cents)
- Decimal places: at most 2
- Validation logic:
  ```go
  amountCents := math.Round(req.Amount * 100)
  if math.Abs(req.Amount*100 - amountCents) > 0.001 {
      return error // more than 2 decimals
  }
  if req.Amount > 2.0 {
      return error // exceeds max
  }
  ```

### 3. No Consecutive Same Recipient
- Query last completed transfer's recipient
- Compare with current transfer's recipient
- Block if they match
- Allow if transferring to a different user

## POST /api/users - Create User

```mermaid
sequenceDiagram
    actor Client
    participant API
    participant UserService
    participant UserRepo
    participant Database

    Client->>API: POST /api/users
    API->>UserService: Create(request)
    
    Note over UserService: Validation: Name length
    alt first_name or last_name > 3 chars
        UserService-->>API: error: name must not exceed 3 characters
        API-->>Client: 400 Bad Request
    end
    
    alt Empty first_name or last_name
        UserService-->>API: error: names are required
        API-->>Client: 400 Bad Request
    end
    
    UserService->>UserRepo: Create(user)
    UserRepo->>Database: INSERT INTO users
    Database-->>UserRepo: userID
    UserRepo-->>UserService: user
    UserService-->>API: user
    API-->>Client: 201 Created
```

## PUT /api/users/:id - Update User

```mermaid
sequenceDiagram
    actor Client
    participant API
    participant UserService
    participant UserRepo
    participant Database

    Client->>API: PUT /api/users/:id
    API->>UserService: Update(id, request)
    
    Note over UserService: Validation: Name length
    alt first_name or last_name > 3 chars
        UserService-->>API: error: name must not exceed 3 characters
        API-->>Client: 400 Bad Request
    end
    
    UserService->>UserRepo: GetByID(id)
    UserRepo->>Database: SELECT * FROM users WHERE id=?
    Database-->>UserRepo: user
    UserRepo-->>UserService: user
    
    alt User not found
        UserService-->>API: error: user not found
        API-->>Client: 404 Not Found
    end
    
    Note over UserService: Merge changes
    UserService->>UserRepo: Update(id, user)
    UserRepo->>Database: UPDATE users SET ...
    Database-->>UserRepo: success
    UserRepo-->>UserService: success
    UserService-->>API: updated user
    API-->>Client: 200 OK
```
