# GitHub Copilot Instructions

> System prompt for GitHub Copilot to follow project coding standards and architectural principles.

## 🎯 Project Overview

**Project Name:** LBK Points - Transfer API  
**Technology Stack:** Go + Fiber (Backend), React + TypeScript (Frontend)  
**Architecture:** Clean Architecture with Repository Pattern  
**Database:** SQLite3  
**API Specification:** OpenAPI 3.1.0

This is a point transfer system where users can:
- Transfer points between users
- View transfer history with pagination
- Check transfer status by idempotency key
- Manage user profiles

## 🏗️ Architecture Principles

### Clean Architecture Layers

```
handlers/ (HTTP Layer)
    ↓
services/ (Business Logic)
    ↓
repositories/ (Data Access)
    ↓
models/ (Data Structures)
```

**MUST:**
- ✅ Always maintain layer separation
- ✅ Dependencies flow downward only (handlers → services → repositories → models)
- ✅ Never skip layers (e.g., handlers calling repositories directly)
- ✅ Use dependency injection for all layers

**MUST NOT:**
- ❌ Never put business logic in handlers
- ❌ Never put database queries in services
- ❌ Never import handlers in services or repositories
- ❌ Never use circular dependencies

## 📝 Coding Standards

### Go Code Style

**File Naming:**
```
✅ repository_user.go
✅ service_transfer.go
✅ handler_user.go
❌ user_repo.go
❌ transferService.go
```

**Function Naming:**
```go
✅ func NewUserService(repo *repositories.UserRepository) *UserService
✅ func (s *UserService) Create(req *models.CreateUserRequest) (*models.User, error)
✅ func (h *UserHandler) GetUsers(c *fiber.Ctx) error

❌ func createNewUserService()  // Not exported when should be
❌ func (s *UserService) create_user()  // Wrong naming convention
```

**Error Handling:**
```go
✅ return fiber.NewError(fiber.StatusBadRequest, "invalid request")
✅ return fmt.Errorf("failed to create user: %w", err)
✅ if err != nil { return err }

❌ panic("error occurred")  // Never panic
❌ log.Fatal(err)  // Don't use in libraries
❌ return errors.New("error")  // Use fmt.Errorf with context
```

### TypeScript/React Code Style

**Component Structure:**
```typescript
✅ function Payment() { ... }  // Functional components
✅ const [state, setState] = useState<Type>(...)

❌ class Payment extends Component { ... }  // Avoid class components
❌ const [state, setState] = useState(null as any)  // Use proper types
```

**Type Safety:**
```typescript
✅ interface User { id: number; name: string }
✅ type TransferStatus = 'pending' | 'completed' | 'failed'
✅ const handleSubmit = (data: FormData): Promise<void> => { ... }

❌ function handleSubmit(data: any) { ... }
❌ let user  // No implicit any
```

## 🔒 Security Guidelines

**MUST:**
- ✅ Always use parameterized queries (prevent SQL injection)
- ✅ Validate all user inputs at handler level
- ✅ Use environment variables for sensitive data
- ✅ Never commit `.env` files or secrets
- ✅ Use transactions for multi-step operations

**Example:**
```go
✅ tx.Exec("INSERT INTO users (name) VALUES (?)", name)
❌ tx.Exec(fmt.Sprintf("INSERT INTO users (name) VALUES ('%s')", name))
```

**MUST NOT:**
- ❌ Never log passwords or tokens
- ❌ Never expose database errors to clients
- ❌ Never trust user input without validation

## 📊 Database Guidelines

**Schema Rules:**
- ✅ Use `snake_case` for column names: `user_id`, `created_at`
- ✅ Always include timestamps: `created_at`, `updated_at`
- ✅ Use foreign keys with proper constraints
- ✅ Create indexes on frequently queried columns
- ✅ Use `INTEGER` for amounts (not FLOAT)

**Query Patterns:**
```go
✅ Use transactions for writes
✅ Use prepared statements
✅ Close rows with defer
✅ Handle sql.ErrNoRows explicitly

❌ Don't use string concatenation for queries
❌ Don't forget to check errors
❌ Don't leave connections open
```

## 🌐 API Design Guidelines

**OpenAPI Compliance:**
- ✅ Follow `swagger.yml` specification exactly
- ✅ Use camelCase for JSON fields: `userId`, `createdAt`
- ✅ Use snake_case for database columns: `user_id`, `created_at`
- ✅ Return proper HTTP status codes

**Response Format:**
```json
✅ { "transfer": { "transferId": 1, ... } }
✅ { "data": [...], "page": 1, "pageSize": 20, "total": 100 }
✅ { "error": "insufficient balance" }

❌ { "TransferId": 1 }  // Wrong case
❌ { "msg": "error" }  // Use "error" key
```

**Status Codes:**
```
✅ 200 OK - Successful GET
✅ 201 Created - Successful POST
✅ 400 Bad Request - Invalid input
✅ 404 Not Found - Resource not found
✅ 409 Conflict - Business rule violation
✅ 422 Unprocessable Entity - Valid format but can't process

❌ 500 Internal Server Error - Only for unexpected errors
❌ 200 OK with error in body - Use proper status codes
```

## 🧪 Testing Guidelines

**Test Structure:**
```go
✅ func TestUserNameValidation(t *testing.T)
✅ func TestTransferAmountValidation(t *testing.T)

❌ func Test1(t *testing.T)  // Not descriptive
❌ func testUserCreation(t *testing.T)  // Must start with Test
```

**Test Patterns:**
```go
✅ Use table-driven tests
✅ Test both success and failure cases
✅ Clean up test data (defer cleanup)
✅ Use test database or in-memory database
✅ Mock external dependencies

❌ Don't test against production database
❌ Don't depend on test execution order
❌ Don't skip cleanup
```

## 📦 Package Management

**Go Modules:**
```bash
✅ go mod tidy  # After adding/removing dependencies
✅ go get github.com/gofiber/fiber/v2
✅ go test ./...

❌ Don't commit vendor/ directory
❌ Don't modify go.mod manually
```

**NPM:**
```bash
✅ npm install <package>
✅ npm run build
✅ npm run dev

❌ Don't commit node_modules/
❌ Don't use npm ci in development
```

## 🚫 Scope Boundaries

### ✅ What Copilot CAN Do

1. **Code Generation:**
   - Generate CRUD operations following repository pattern
   - Create API handlers with proper validation
   - Write unit tests for services and repositories
   - Generate OpenAPI/Swagger documentation

2. **Refactoring:**
   - Split large functions into smaller ones
   - Extract common logic into utilities
   - Improve error handling
   - Add missing error checks

3. **Documentation:**
   - Add function comments
   - Generate README sections
   - Create API documentation
   - Write inline code comments

4. **Debugging Assistance:**
   - Identify potential bugs
   - Suggest fixes for compilation errors
   - Explain error messages
   - Recommend better error handling

### ❌ What Copilot CANNOT Do

1. **Architecture Changes:**
   - ❌ Don't change the layer structure without discussion
   - ❌ Don't introduce new architectural patterns arbitrarily
   - ❌ Don't modify database schema without migration plan
   - ❌ Don't change API contracts that break clients

2. **Security Decisions:**
   - ❌ Don't generate authentication/authorization logic without requirements
   - ❌ Don't decide on encryption methods
   - ❌ Don't implement rate limiting without specifications
   - ❌ Don't add CORS rules without approval

3. **Business Logic:**
   - ❌ Don't invent business rules (e.g., transfer limits, validation rules)
   - ❌ Don't modify calculation logic without confirmation
   - ❌ Don't change status flow (pending → processing → completed)
   - ❌ Don't add new transfer states arbitrarily

4. **External Integrations:**
   - ❌ Don't add external API calls without discussion
   - ❌ Don't integrate third-party services
   - ❌ Don't add payment gateways
   - ❌ Don't implement notification systems

5. **Production Changes:**
   - ❌ Don't modify production configurations
   - ❌ Don't generate migration scripts that run automatically
   - ❌ Don't change environment variables in production
   - ❌ Don't modify CI/CD pipelines without review

## 🎨 Code Formatting

**Go:**
```bash
✅ Use gofmt (automatic)
✅ Use goimports for imports
✅ Follow https://go.dev/doc/effective_go
```

**TypeScript/React:**
```bash
✅ Use ESLint configuration
✅ Use Prettier for formatting
✅ Follow React best practices
```

## 📚 Documentation Requirements

**Function Comments:**
```go
✅ // NewUserService creates a new user service with the given repository
func NewUserService(repo *repositories.UserRepository) *UserService

✅ // Create validates the input and creates a new user in the database.
// Returns an error if validation fails or database operation fails.
func (s *UserService) Create(req *models.CreateUserRequest) (*models.User, error)

❌ func Create(req *models.CreateUserRequest) (*models.User, error)  // No comment
```

**Complex Logic:**
```go
✅ // Check if user has sufficient balance before transfer
// Balance must be at least equal to transfer amount
if fromBalance < req.Amount {
    return nil, errors.New("insufficient balance")
}

✅ // Prevent consecutive transfers to the same recipient
// This is a business rule to prevent spam
lastRecipient, err := s.userRepo.GetLastTransferRecipient(req.FromUserID)
```

## 🔄 Git Commit Guidelines

**Commit Message Format:**
```
✅ feat: Add transfer history pagination
✅ fix: Resolve balance update race condition
✅ docs: Update API documentation for transfers
✅ refactor: Extract validation logic to separate function
✅ test: Add unit tests for transfer service

❌ Update code
❌ Fix bug
❌ WIP
```

**Types:**
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation only
- `refactor`: Code change that neither fixes a bug nor adds a feature
- `test`: Adding missing tests
- `chore`: Changes to build process or auxiliary tools

## 🚀 Performance Guidelines

**MUST:**
- ✅ Use indexes for frequently queried columns
- ✅ Paginate large result sets
- ✅ Use connection pooling
- ✅ Close database connections properly
- ✅ Use transactions efficiently

**MUST NOT:**
- ❌ N+1 query problems
- ❌ Loading entire tables into memory
- ❌ Unnecessary database queries in loops
- ❌ Missing indexes on foreign keys

## 📖 References

- **Project Style Guide:** `/STYLE_GUIDE.md`
- **API Specification:** `/backend/swagger.yml`
- **Database Schema:** `/backend/database.md`
- **MCP Setup:** `/MCP_SETUP.md`
- **Workshop Reference:** https://github.com/mikelopster/kbtg-ai-workshop-oct

## 🎯 Priority Rules

When generating code, follow this priority order:

1. **Security First** - Always validate input, use parameterized queries
2. **Follow Architecture** - Respect layer boundaries
3. **Match Specifications** - Follow swagger.yml exactly
4. **Code Quality** - Readable, maintainable, well-documented
5. **Performance** - Efficient queries, proper indexing
6. **Error Handling** - Comprehensive error checking and reporting

---

**Last Updated:** October 17, 2025  
**Version:** 1.0.0  
**Maintainer:** Development Team
