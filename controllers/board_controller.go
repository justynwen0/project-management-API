package controllers

import (
	"math"
	"strconv"

	"github.com/Filbertfelix888/project-management-API/models"
	"github.com/Filbertfelix888/project-management-API/services"
	"github.com/Filbertfelix888/project-management-API/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

type BoardController struct {
	services services.BoardService
}

func NewBoardController(s services.BoardService) *BoardController {
	return &BoardController{services: s}
}

func (c *BoardController) CreateBoard(ctx *fiber.Ctx) error {
	var userID uuid.UUID
	var err error

	board := new(models.Board)
	user := ctx.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	if err := ctx.BodyParser(board); err != nil {
		return utils.BadRequest(ctx, "Gagal membaca request", err.Error())
	}

	userID, err = uuid.Parse(claims["pub_id"].(string))
	if err != nil {
		return utils.BadRequest(ctx, "ID pengguna tidak valid", err.Error())
	}
	board.OwnerPublicID = userID

	if err := c.services.Create(board); err != nil {
		return utils.BadRequest(ctx, "Gagal menyimpan data", err.Error())
	}

	return utils.Success(ctx, "Board berhasil dibuat", board)
}

func (c *BoardController) UpdateBoard(ctx *fiber.Ctx) error {
	publicID := ctx.Params("id")
	board := new(models.Board)
	if err := ctx.BodyParser(board); err != nil {
		return utils.BadRequest(ctx, "Gagal parsing data", err.Error())
	}
	if _, err := uuid.Parse(publicID); err != nil {
		return utils.BadRequest(ctx, "ID board tidak valid", err.Error())
	}
	exitingBoard, err := c.services.GetByPublicID(publicID)
	if err != nil {
		return utils.NotFound(ctx, "Board tidak ditemukan", err.Error())
	}
	board.InternalID = exitingBoard.InternalID
	board.PublicID = exitingBoard.PublicID
	board.OwnerID = exitingBoard.OwnerID
	board.OwnerPublicID = exitingBoard.OwnerPublicID
	board.CreatedAt = exitingBoard.CreatedAt

	if err := c.services.Update(board); err != nil {
		return utils.BadRequest(ctx, "Gagal update board", err.Error())
	}
	return utils.Success(ctx, "Board berhasil diperbarui", board)
}

func (c *BoardController) AddBoardMember(ctx *fiber.Ctx) error {
	publicID := ctx.Params("id")

	var userIDs []string
	if err := ctx.BodyParser(&userIDs); err != nil {
		return utils.BadRequest(ctx, "Gagal parsing data", err.Error())
	}

	if err := c.services.AddMembers(publicID, userIDs); err != nil {
		return utils.BadRequest(ctx, "Gagal menambahkan members", err.Error())
	}
	return utils.Success(ctx, "Members berhasil ditambahkan", nil)
}

func (c *BoardController) RemoveBoardMember(ctx *fiber.Ctx) error {
	publicID := ctx.Params("id")

	var userIDs []string
	if err := ctx.BodyParser(&userIDs); err != nil {
		return utils.BadRequest(ctx, "Gagal parsing data", err.Error())
	}

	if err := c.services.RemoveMembers(publicID, userIDs); err != nil {
		return utils.BadRequest(ctx, "Gagal menghapus members", err.Error())
	}
	return utils.Success(ctx, "Members berhasil dihapus", nil)
}

func (c *BoardController) GetMyBoardPaginate(ctx *fiber.Ctx) error {
	user := ctx.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := claims["pub_id"].(string)

	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	limit, _ := strconv.Atoi(ctx.Query("limit", "10"))
	offset := (page - 1) * limit

	filter := ctx.Query("filter", "")
	sort := ctx.Query("sort", "created_at desc")

	board, total, err := c.services.GetAllByUserPaginate(userID, filter, sort, limit, offset)
	if err != nil {
		return utils.InternalServerError(ctx, "Gagal mengambil data board", err.Error())
	}

	meta := utils.PaginationMeta{
		Page:      page,
		Limit:     limit,
		Total:     int(total),
		TotalPage: int(math.Ceil(float64(total) / float64(limit))),
		Filter:    filter,
		Sort:      sort,
	}
	return utils.SuccessPagination(ctx, "Data board berhasil diambil", board, meta)
}

func (c *BoardController) GetBoardById(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	board, err := c.services.GetByPublicID(id)
	if err != nil {
		return utils.NotFound(ctx, "Board tidak ditemukan", err.Error())
	}

	var boardResp models.Board
	err = copier.Copy(&boardResp, &board)

	if err != nil {
		return utils.BadRequest(ctx, "Internal Server Error", err.Error())
	}

	return utils.Success(ctx, "Board berhasil diambil", boardResp)
}

func (c *BoardController) GetBoardMembers(ctx *fiber.Ctx) error {
	boardPublicID := ctx.Params("board_id")

	if _, err := uuid.Parse(boardPublicID); err != nil {
		return utils.BadRequest(ctx, "ID Board tidak valid", err.Error())
	}

	boardMembers, err := c.services.GetMembersByBoardID(boardPublicID)
	if err != nil {
		return utils.NotFound(ctx, "Members tidak ditemukan untuk board ini", err.Error())
	}
	return utils.Success(ctx, "Members board berhasil diambil", boardMembers)
}