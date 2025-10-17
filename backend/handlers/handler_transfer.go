package handlers

import (
	"backend/models"
	"backend/services"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type TransferHandler struct {
	service *services.TransferService
}

func NewTransferHandler(service *services.TransferService) *TransferHandler {
	return &TransferHandler{service: service}
}

func (h *TransferHandler) CreateTransfer(c *fiber.Ctx) error {
	var req models.CreateTransferRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	transfer, err := h.service.CreateTransfer(&req)
	if err != nil {
		return err
	}

	// Set Idempotency-Key header
	c.Set("Idempotency-Key", transfer.IdemKey)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"transfer": transfer,
	})
}

func (h *TransferHandler) GetTransfer(c *fiber.Ctx) error {
	idemKey := c.Params("id")

	transfer, err := h.service.GetByIdemKey(idemKey)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	return c.JSON(fiber.Map{
		"transfer": transfer,
	})
}

func (h *TransferHandler) ListTransfers(c *fiber.Ctx) error {
	userIDStr := c.Query("userId")
	if userIDStr == "" {
		return fiber.NewError(fiber.StatusBadRequest, "userId query parameter is required")
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid userId")
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize", "20"))

	result, err := h.service.GetByUserID(userID, page, pageSize)
	if err != nil {
		return err
	}

	return c.JSON(result)
}
