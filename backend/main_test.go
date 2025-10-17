package main

import (
	"backend/handlers"
	"backend/models"
	"backend/repositories"
	"backend/services"
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func setupTestApp(t *testing.T) (*fiber.App, *sql.DB) {
	db, err := InitDB(":memory:")
	if err != nil {
		t.Fatalf("Failed to init test DB: %v", err)
	}

	if err := Migrate(db); err != nil {
		t.Fatalf("Failed to migrate test DB: %v", err)
	}

	userRepo := repositories.NewUserRepository(db)
	transferRepo := repositories.NewTransferRepository(db)
	ledgerRepo := repositories.NewLedgerRepository(db)

	userService := services.NewUserService(userRepo)
	transferService := services.NewTransferService(transferRepo, ledgerRepo, userRepo)

	userHandler := handlers.NewUserHandler(userService)
	transferHandler := handlers.NewTransferHandler(transferService)

	app := fiber.New(fiber.Config{
		ErrorHandler: ErrorHandler,
	})

	api := app.Group("/api")
	users := api.Group("/users")
	users.Get("/", userHandler.GetUsers)
	users.Get("/:id", userHandler.GetUser)
	users.Post("/", userHandler.CreateUser)
	users.Put("/:id", userHandler.UpdateUser)
	users.Delete("/:id", userHandler.DeleteUser)

	transfers := api.Group("/transfers")
	transfers.Post("/", transferHandler.CreateTransfer)
	transfers.Get("/:id", transferHandler.GetTransfer)

	return app, db
}

// Test Case 1: Names must not exceed 3 characters
func TestUserNameValidation(t *testing.T) {
	app, db := setupTestApp(t)
	defer db.Close()

	tests := []struct {
		name      string
		firstName string
		lastName  string
		wantError bool
	}{
		{"Valid names", "Tom", "Lee", false},
		{"First name too long", "Thomas", "Lee", true},
		{"Last name too long", "Tom", "Johnson", true},
		{"Both names too long", "Thomas", "Johnson", true},
		{"Empty names", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := models.CreateUserRequest{
				FirstName: tt.firstName,
				LastName:  tt.lastName,
				Email:     "test@example.com",
			}
			jsonBody, _ := json.Marshal(body)

			req := httptest.NewRequest("POST", "/api/users", bytes.NewReader(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			if err != nil {
				t.Fatal(err)
			}

			if tt.wantError {
				if resp.StatusCode == 201 {
					t.Errorf("Expected error but got status 201")
				}
			} else {
				if resp.StatusCode != 201 {
					t.Errorf("Expected status 201 but got %d", resp.StatusCode)
				}
			}
		})
	}
}

// Test Case 2: Transfer amount limits
func TestTransferAmountValidation(t *testing.T) {
	app, db := setupTestApp(t)
	defer db.Close()

	tests := []struct {
		name      string
		amount    int64
		wantError bool
		errorMsg  string
	}{
		{"Valid amount", 150, false, ""},
		{"Large amount", 10000, false, ""},
		{"Negative amount", -100, true, "amount must be greater than 0"},
		{"Zero amount", 0, true, "amount must be greater than 0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create fresh users for each test to avoid balance issues
			user1ID := createTestUserWithBalance(t, db, "Bob", "Cat", 10000)
			user2ID := createTestUserWithBalance(t, db, "Ann", "Doe", 0)

			body := models.CreateTransferRequest{
				FromUserID: user1ID,
				ToUserID:   user2ID,
				Amount:     tt.amount,
			}
			jsonBody, _ := json.Marshal(body)

			req := httptest.NewRequest("POST", "/api/transfers", bytes.NewReader(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			if err != nil {
				t.Fatal(err)
			}

			if tt.wantError {
				if resp.StatusCode == 201 {
					t.Errorf("Expected error but got status 201")
				}
			} else {
				if resp.StatusCode != 201 {
					var errResp map[string]interface{}
					json.NewDecoder(resp.Body).Decode(&errResp)
					t.Errorf("Expected status 201 but got %d, error: %v", resp.StatusCode, errResp)
				}
			}
		})
	}
}

// Test Case 3: Cannot transfer to same recipient consecutively
func TestNoConsecutiveSameRecipient(t *testing.T) {
	app, db := setupTestApp(t)
	defer db.Close()

	// Create test users with sufficient balance
	userA := createTestUserWithBalance(t, db, "Aki", "Tan", 50000) // 500.00 balance (enough for 5 transfers)
	userB := createTestUserWithBalance(t, db, "Ben", "Kim", 0)
	userC := createTestUserWithBalance(t, db, "Cat", "Lin", 0)

	// First transfer: A -> B (should succeed)
	transfer1 := models.CreateTransferRequest{
		FromUserID: userA,
		ToUserID:   userB,
		Amount:     100,
	}
	jsonBody1, _ := json.Marshal(transfer1)
	req1 := httptest.NewRequest("POST", "/api/transfers", bytes.NewReader(jsonBody1))
	req1.Header.Set("Content-Type", "application/json")
	resp1, _ := app.Test(req1)

	if resp1.StatusCode != 201 {
		t.Fatalf("First transfer failed with status %d", resp1.StatusCode)
	}

	// Second transfer: A -> B again (should fail)
	transfer2 := models.CreateTransferRequest{
		FromUserID: userA,
		ToUserID:   userB,
		Amount:     100,
	}
	jsonBody2, _ := json.Marshal(transfer2)
	req2 := httptest.NewRequest("POST", "/api/transfers", bytes.NewReader(jsonBody2))
	req2.Header.Set("Content-Type", "application/json")
	resp2, _ := app.Test(req2)

	if resp2.StatusCode == 201 {
		t.Errorf("Second transfer to same recipient should have failed but got status 201")
	}

	// Third transfer: A -> C (should succeed)
	transfer3 := models.CreateTransferRequest{
		FromUserID: userA,
		ToUserID:   userC,
		Amount:     100,
	}
	jsonBody3, _ := json.Marshal(transfer3)
	req3 := httptest.NewRequest("POST", "/api/transfers", bytes.NewReader(jsonBody3))
	req3.Header.Set("Content-Type", "application/json")
	resp3, _ := app.Test(req3)

	if resp3.StatusCode != 201 {
		t.Errorf("Transfer to different recipient should succeed but got status %d", resp3.StatusCode)
	}

	// Fourth transfer: A -> B (should succeed now since last was to C)
	transfer4 := models.CreateTransferRequest{
		FromUserID: userA,
		ToUserID:   userB,
		Amount:     100,
	}
	jsonBody4, _ := json.Marshal(transfer4)
	req4 := httptest.NewRequest("POST", "/api/transfers", bytes.NewReader(jsonBody4))
	req4.Header.Set("Content-Type", "application/json")
	resp4, _ := app.Test(req4)

	if resp4.StatusCode != 201 {
		var errResp map[string]interface{}
		json.NewDecoder(resp4.Body).Decode(&errResp)
		t.Errorf("Transfer to B should succeed after transferring to C, but got status %d, error: %v", resp4.StatusCode, errResp)
	}
}

func createTestUserWithBalance(t *testing.T, db *sql.DB, firstName, lastName string, balance int64) int64 {
	now := models.Now()
	result, err := db.Exec(`
		INSERT INTO users (first_name, last_name, email, phone, avatar_url, bio, points_balance, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, firstName, lastName, "", "", "", "", balance, now, now)

	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	id, _ := result.LastInsertId()
	return id
}
