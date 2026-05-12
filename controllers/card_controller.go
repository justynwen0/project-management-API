package controllers

import (
	"time"

	"github.com/Filbertfelix888/project-management-API/models"
	"github.com/Filbertfelix888/project-management-API/services"
	"github.com/Filbertfelix888/project-management-API/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type CardController struct {
	service services.CardService
}

func NewCardController(s services.CardService) *CardController {
	return &CardController{service: s}
}

func (c *CardController) CreateCard(ctx *fiber.Ctx) error {
	type CreateCardRequest struct {
		ListPublicID string    `json:"list_id"`
		Title        string    `json:"title"`
		Description  string    `json:"description"`
		DueDate      time.Time `json:"due_date"`
		Position     int       `json:"position"`
	}
	var req CreateCardRequest
	if err := ctx.BodyParser(&req); err != nil {
		return utils.BadRequest(ctx, "Gagal mengambil data", err.Error())
	}

	card := &models.Card{
		Title:       req.Title,
		Description: req.Description,
		DueDate:     &req.DueDate,
		Position:    req.Position,
	}

	if err := c.service.Create(card, req.ListPublicID); err != nil {
		return utils.InternalServerError(ctx, "Gagal membuat card", err.Error())
	}
	return utils.Success(ctx, "Card berhasil dibuat", card)
}

func (c *CardController) UpdateCard(ctx *fiber.Ctx) error {
	publicID := ctx.Params("id")

	type updateCardRequest struct {
		ListPublicID string `json:"list_id"`
		Title string `json:"title"`
		Description string `json:"description"`
		DueDate *time.Time `json:"due_date"`
		Position int `json:"position"`
	}

	var req updateCardRequest
	if err := ctx.BodyParser(&req); err != nil {
		return utils.BadRequest(ctx, "Gagal parsing data", err.Error())
	}

	if _, err := uuid.Parse(publicID); err != nil {
		return utils.BadRequest(ctx, "ID tidak valid", err.Error())
	}

	card := &models.Card{
		Title: req.Title,
		Description: req.Description,
		DueDate: req.DueDate,
		Position: req.Position,
		PublicID: uuid.MustParse(publicID),
	}
	if err := c.service.Update(card, req.ListPublicID); err != nil {
		return utils.InternalServerError(ctx, "Gagal update data", err.Error())
	}
	return  utils.Success(ctx, "Card berhasil diperbaharui", card)
}

func (c *CardController) DeleteCard(ctx *fiber.Ctx) error {
	publicID := ctx.Params("id")

	if _, err := uuid.Parse(publicID); err != nil {
		return utils.BadRequest(ctx, "ID tidak valid", err.Error())
	}

	card, err := c.service.GetByPublicID(publicID)
	if err != nil {
		return utils.NotFound(ctx, "Card tidak ditemukan", err.Error())
	}

	if err := c.service.Delete(uint(card.InternalID)); err != nil {
		return utils.BadRequest(ctx, "Gagal menghapus data", err.Error())
	}
	return utils.Success(ctx, "Card berhasil dihapus", publicID)
}

func (c *CardController) GetCardDetail(ctx *fiber.Ctx) error {
	cardPublicID := ctx.Params("id")

	card, err := c.service.GetByPublicID(cardPublicID)
	if err != nil {
		return utils.InternalServerError(ctx, "Error saat mengambil data", err.Error())
	}
	if card == nil {
		return utils.NotFound(ctx, "Card tidak ditemukan", err.Error())
	}
	return utils.Success(ctx, "Data berhasil diambil", card)
}

func (c *CardController) GetCardOnList(ctx *fiber.Ctx) error {
	listPublicID := ctx.Params("id")

	if _, err := uuid.Parse(listPublicID); err != nil {
		return utils.BadRequest(ctx, "ID daftar tidak valid", err.Error())
	}

	cards, _ := c.service.GetByListID(listPublicID)
	// if err != nil {
	// 	return utils.NotFound(ctx, "Cards not found for the list", err.Error())
	// }
	return utils.Success(ctx, "Cards berhasil diambil", cards)
}

func (c *CardController) AddCardAssignees(ctx *fiber.Ctx) error {
    cardID := ctx.Params("id")

    type requestBody struct {
        UserID  []string `json:"user_id"`
        UserIDs []string `json:"user_ids"`
    }

    var req requestBody
    if err := ctx.BodyParser(&req); err != nil {
        return utils.BadRequest(ctx, "Gagal parsing data", err.Error())
    }

    userIDs := req.UserID
    if len(userIDs) == 0 {
        userIDs = req.UserIDs
    }

    if len(userIDs) == 0 {
        return utils.BadRequest(ctx, "Data assignee tidak ditemukan", "user_id atau user_ids wajib diisi")
    }
	// Pastikan memanggil fungsi yang menerima PublicID (string)
    
	if err := c.service.AddAssigneesByPublicID(cardID, userIDs); err != nil {
        return utils.InternalServerError(ctx, "Gagal menambahkan assignee", err.Error())
    }
    return utils.Success(ctx, "Assignee berhasil ditambahkan", nil)
}	

func (c *CardController) RemoveCardAssignees(ctx *fiber.Ctx) error {
    cardID := ctx.Params("id")

	type requestBody struct {
		UserID  []string `json:"user_id"`
		UserIDs []string `json:"user_ids"`
	}
	
    var req requestBody
    if err := ctx.BodyParser(&req); err != nil {
        return utils.BadRequest(ctx, "Gagal parsing data", err.Error())
    }

    userIDs := req.UserID
    if len(userIDs) == 0 {
        userIDs = req.UserIDs
    }

    if err := c.service.RemoveAssignees(cardID, userIDs); err != nil {
        return utils.InternalServerError(ctx, "Gagal menghapus assignee", err.Error())
    }
    return utils.Success(ctx, "Assignee berhasil dihapus", nil)
}