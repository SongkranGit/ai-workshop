package main

import (
	"backend/handlers"
	"backend/repositories"
	"backend/services"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// Initialize database
	db, err := InitDB("./data.db")
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// Run migrations
	if err := Migrate(db); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	transferRepo := repositories.NewTransferRepository(db)
	ledgerRepo := repositories.NewLedgerRepository(db)

	// Initialize services
	userService := services.NewUserService(userRepo)
	transferService := services.NewTransferService(transferRepo, ledgerRepo, userRepo)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService)
	transferHandler := handlers.NewTransferHandler(transferService)

	// Setup Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: ErrorHandler,
	})

	// Middleware
	app.Use(cors.New())
	app.Use(logger.New())

	// Swagger UI
	SetupSwagger(app)

	// Routes
	api := app.Group("/api")

	// User routes
	users := api.Group("/users")
	users.Get("/", userHandler.GetUsers)
	users.Get("/:id", userHandler.GetUser)
	users.Post("/", userHandler.CreateUser)
	users.Put("/:id", userHandler.UpdateUser)
	users.Delete("/:id", userHandler.DeleteUser)

	// Transfer routes
	transfers := api.Group("/transfers")
	transfers.Post("/", transferHandler.CreateTransfer)
	transfers.Get("/", transferHandler.ListTransfers)
	transfers.Get("/:id", transferHandler.GetTransfer)

	log.Println("Server starting on :3000")
	log.Println("Swagger UI available at: http://localhost:3000/swagger")
	if err := app.Listen(":3000"); err != nil {
		log.Fatal(err)
	}
}
