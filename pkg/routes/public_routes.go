package routes

import (
	"github.com/Figbase/api/app/controllers"
	"github.com/gofiber/fiber/v2"
)

// PublicRoutes func for describe group of public routes.
func PublicRoutes(a *fiber.App) {
	// Create routes group.
	route := a.Group("/api/v1")

	// Routes for GET method:
	route.Get("/", controllers.Home) // get home route
	// route.Get("/book/:id", controllers.GetBook) // get one book by ID

	// Routes for POST method:
	route.Post("/auth/signup", controllers.UserSignUp) // register a new user
	route.Post("/auth/signin", controllers.UserSignIn) // auth, return Access & Refresh tokens
}
