package routes

import (
	"module/controllers"
	"module/middleware"


	"github.com/gofiber/fiber/v2"
)

func UserRoutes(app *fiber.App){
	app.Post("/User/Register", controllers.Register)



	app.Get("/Users", middleware.JWTProtected(), controllers.GetUsers)
	app.Get("/User/Delete/:userId",middleware.JWTProtected(),controllers.DeleteUser)
}