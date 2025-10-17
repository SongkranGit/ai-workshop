# Database Schema

```mermaid
erDiagram
    users ||--o{ transfers : "sends (from)"
    users ||--o{ transfers : "receives (to)"
    users ||--o{ point_ledger : "has"
    transfers ||--o{ point_ledger : "generates"

    users {
        INTEGER id PK
        TEXT first_name "NOT NULL, max 3 chars"
        TEXT last_name "NOT NULL, max 3 chars"
        TEXT email
        TEXT phone
        TEXT avatar_url
        TEXT bio
        INTEGER points_balance "NOT NULL, default 0"
        TEXT created_at "NOT NULL"
        TEXT updated_at "NOT NULL"
    }

    transfers {
        INTEGER id PK
        INTEGER from_user_id FK "NOT NULL"
        INTEGER to_user_id FK "NOT NULL"
        INTEGER amount "NOT NULL, in cents, >0, <=200"
        TEXT status "NOT NULL, enum"
        TEXT note
        TEXT idempotency_key "NOT NULL, UNIQUE"
        TEXT created_at "NOT NULL"
        TEXT updated_at "NOT NULL"
        TEXT completed_at
        TEXT fail_reason
    }

    point_ledger {
        INTEGER id PK
        INTEGER user_id FK "NOT NULL"
        INTEGER change "NOT NULL"
        INTEGER balance_after "NOT NULL"
        TEXT event_type "NOT NULL, enum"
        INTEGER transfer_id FK
        TEXT reference
        TEXT metadata "JSON"
        TEXT created_at "NOT NULL"
    }
```

## Business Rules

### Users
- `first_name` and `last_name` must not exceed 3 characters
- `points_balance` tracks user's current point balance (in cents)

### Transfers
- Maximum transfer amount: 2.00 (200 cents)
- Amount must have at most 2 decimal places
- Cannot transfer to the same recipient as the last completed transfer
- Status values: pending, processing, completed, failed, cancelled, reversed
- `idempotency_key` ensures duplicate prevention

### Point Ledger
- Append-only transaction log
- Event types: transfer_out, transfer_in, adjust, earn, redeem
- Each transfer creates two ledger entries (one for sender, one for receiver)
- `balance_after` provides audit trail of balance changes
