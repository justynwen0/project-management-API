package main

import (
	"log"

	"github.com/Filbertfelix888/project-management-API/config"
	"github.com/Filbertfelix888/project-management-API/controllers"
	"github.com/Filbertfelix888/project-management-API/database/seed"
	"github.com/Filbertfelix888/project-management-API/repositories"
	"github.com/Filbertfelix888/project-management-API/routes"
	"github.com/Filbertfelix888/project-management-API/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	config.LoadEnv()
	config.ConnectDB()

	seed.SeedAdmin()
	app := fiber.New()

	app.Use(cors.New(cors.Config{
        AllowOrigins:     "http://localhost:5173",
        AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
        AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
        AllowCredentials: true,
    }))

	// user
	userRepo := repositories.NewUserRepository()
	userService := services.NewUserService(userRepo)
	userController := controllers.NewUserController(userService)

	// board
	boardRepo := repositories.NewBoardRepository()
	boardMemberRepo := repositories.NewBoardMemberRepository()
	boardService := services.NewBoardService(boardRepo, userRepo, boardMemberRepo)
	boardController := controllers.NewBoardController(boardService)

	// list
	listPosRepo := repositories.NewListPositionRepository()
	listRepo := repositories.NewListRepositories()
	listService := services.NewListService(listRepo, boardRepo, listPosRepo)
	listController := controllers.NewListController(listService)

	// card
	cardRepo := repositories.NewCardRepository()
	cardService := services.NewCardService(cardRepo, listRepo, userRepo)
	cardController := controllers.NewCardController(cardService)

	routes.Setup(app, userController, boardController, listController, cardController)

	port := config.AppConfig.AppPort
	log.Println("Server is running on port :", port)
	log.Fatal(app.Listen(":" + port))
}
