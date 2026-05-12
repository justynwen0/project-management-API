package controllers

import (
	"github.com/Filbertfelix888/project-management-API/models"
	"github.com/Filbertfelix888/project-management-API/services"
	"github.com/Filbertfelix888/project-management-API/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ListController struct {
	service services.ListService
}

func NewListController(s services.ListService) *ListController {
	return &ListController{service: s}
}

func (c *ListController) CreateList(ctx *fiber.Ctx) error {
	list := new(models.List)
	if err := ctx.BodyParser(list); err != nil {
		return utils.BadRequest(ctx, "Gagal membaca request", err.Error())
	}
	if err := c.service.Create(list); err != nil {
		return utils.BadRequest(ctx, "Gagal membuat list", err.Error())
	}
	return utils.Success(ctx, "List berhasil dibuat", list)
}

func (c *ListController) UpdateList(ctx *fiber.Ctx) error {
	publicID := ctx.Params("id")
	list := new(models.List)

	if err := ctx.BodyParser(list); err != nil {
		return utils.BadRequest(ctx, "Gagal parsing data", err.Error())
	}

	if _, err := uuid.Parse(publicID); err != nil {
		return utils.BadRequest(ctx, "ID tidak valid", err.Error())
	}

	exitingList, err := c.service.GetByPublicID(publicID)
	if err != nil {
		return utils.NotFound(ctx, "List tidak ditemukan", err.Error())
	}

	list.InternalID = exitingList.InternalID
	list.PublicID = exitingList.PublicID

	if err := c.service.Update(list); err != nil {
		return utils.BadRequest(ctx, "Gagal memperbarui list", err.Error())
	}
	
	updatedList, err := c.service.GetByPublicID(publicID)
	if err != nil {
		return utils.NotFound(ctx, "List tidak ditemukan", err.Error())
	}
	return utils.Success(ctx, "List berhasil diperbarui", updatedList)
}

func (c *ListController) GetListOnBoard(ctx *fiber.Ctx) error {
	boardPublicID := ctx.Params("board_id")
	if _, err := uuid.Parse(boardPublicID); err != nil {
		return utils.BadRequest(ctx, "ID tidak valid", err.Error())
	}
	lists, err := c.service.GetByBoardID(boardPublicID)
	if err != nil {
		return utils.NotFound(ctx, "List tidak ditemukan", err.Error())
	}
	return utils.Success(ctx, "Data berhasil diambil", lists)
}

func (c *ListController) DeleteList(ctx *fiber.Ctx) error {
	publicID := ctx.Params("id")
	if _, err := uuid.Parse(publicID); err != nil {
		return utils.BadRequest(ctx, "ID tidak valid", err.Error())
	}

	list, err := c.service.GetByPublicID(publicID)
	if err != nil {
		return utils.NotFound(ctx, "List tidak ditemukan", err.Error())
	}

	if err := c.service.Delete(uint(list.InternalID)); err != nil {
		return utils.InternalServerError(ctx, "Gagal menghapus list", err.Error())
	}

	return utils.Success(ctx, "List berhasil dihapus", list)
}