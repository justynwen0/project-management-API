package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/justynwen0/project-management-API/services"
)

type DashboardController struct {
	dashboardService services.DashboardService
}

func NewDashboardController(
	dashboardService services.DashboardService,
) *DashboardController {
	return &DashboardController{
		dashboardService: dashboardService,
	}
}

func (c *DashboardController) GetWorkload(ctx *fiber.Ctx) error {
	data, err := c.dashboardService.GetTaskCountByAssignee()

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return ctx.JSON(data)
}

func (c *DashboardController) GetTaskPercentage(ctx *fiber.Ctx) error {
	data, err := c.dashboardService.GetTaskPercentage()

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return ctx.JSON(data)
}
