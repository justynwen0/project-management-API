package routes

import (
	"log"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/joho/godotenv"
	"github.com/justynwen0/project-management-API/config"
	"github.com/justynwen0/project-management-API/controllers"
	"github.com/justynwen0/project-management-API/utils"
)

func Setup(app *fiber.App,
	uc *controllers.UserController,
	bc *controllers.BoardController,
	lc *controllers.ListController,
	cc *controllers.CardController,
	dc *controllers.DashboardController) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	app.Post("/v1/auth/register", uc.Register)
	app.Post("/v1/auth/login", uc.Login)

	// JWT protected routes
	api := app.Group("/api/v1", jwtware.New(jwtware.Config{
		SigningKey: []byte(config.AppConfig.JWTSecret),
		ContextKey: "user",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return utils.Unauthorized(c, "Error unauthorized", err.Error())
		},
	}))

	//user
	userGroup := api.Group("/users")
	userGroup.Get("/:page", uc.GetUserPagination) // /api/v1/users/page?filter=&sort=-id&page=1&limit=20
	userGroup.Get("/:id", uc.GetUser)             //  /api/v1/users/id
	userGroup.Put("/:id", uc.UpdateUser)          // /api/v1/users/id
	userGroup.Delete("/:id", uc.DeleteUser)       // /api/v1/users/id

	//board
	boardGroup := api.Group("/boards")
	boardGroup.Post("/", bc.CreateBoard)                    // /api/v1/boards
	boardGroup.Put("/:id", bc.UpdateBoard)                  // /api/v1/boards/id
	boardGroup.Delete("/:id", bc.DeleteBoard)                // /api/v1/boards/id
	boardGroup.Post("/:id/members", bc.AddBoardMember)      // /api/v1/boards/id/members
	boardGroup.Delete("/:id/members", bc.RemoveBoardMember) // /api/v1/boards/id/members
	boardGroup.Get("/my", bc.GetMyBoardPaginate)            // /api/v1/boards/my
	boardGroup.Get("/:board_id/lists", lc.GetListOnBoard)   // /api/v1/boards/board_id
	boardGroup.Get("/:id", bc.GetBoardById)                 // /api/v1/boards/id
	boardGroup.Get("/:board_id/members", bc.GetBoardMembers)

	//list
	listGroup := api.Group("/lists")
	listGroup.Get("/:id/cards", cc.GetCardOnList) // /api/v1/lists/id/cardsp
	listGroup.Post("/", lc.CreateList)            // /api/v1/lists
	listGroup.Put("/:id", lc.UpdateList)          // /api/v1/lists/id
	listGroup.Delete("/:id", lc.DeleteList)       // /api/v1/lists/id

	//card
	cardGroup := api.Group("/cards")
	cardGroup.Post("/", cc.CreateCard)
	cardGroup.Put("/:id", cc.UpdateCard)
	cardGroup.Delete("/:id", cc.DeleteCard)
	cardGroup.Get("/:id", cc.GetCardDetail)
	cardGroup.Post("/:id/assignees", cc.AddCardAssignees)
	cardGroup.Delete("/:id/assignees", cc.RemoveCardAssignees)

	//dashboard
	dashboardGroup := api.Group("/dashboard")
	dashboardGroup.Get("/workload", dc.GetWorkload)
	dashboardGroup.Get("/task-percentage", dc.GetTaskPercentage)

	app.Get("/test-workload", dc.GetWorkload)
	app.Get("/test-task-percentage", dc.GetTaskPercentage)

}
