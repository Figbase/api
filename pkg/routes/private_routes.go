package routes

import (
	"github.com/Figbase/api/app/controllers"
	"github.com/Figbase/api/pkg/middleware"
	"github.com/gofiber/fiber/v2"
)

const (
	apiKey = "my-super-secret-key"
)

// PrivateRoutes func for describe group of private routes.
func PrivateRoutes(a *fiber.App) {
	// Create routes group.
	route := a.Group("/api/v1")

	// Routes for POST method:
	// route.Post("/book", middleware.JWTProtected(), controllers.CreateBook)           // create a new book
	route.Post("/auth/signout", middleware.JWTProtected(), controllers.UserSignOut) // de-authorization user
	route.Post("/token/renew", middleware.JWTProtected(), controllers.RenewTokens)  // renew Access & Refresh tokens

	// route.Post("/api-key", middleware.AuthMiddleware(apiKey), controllers.Home) // renew Access & Refresh tokens

	// Routes for PUT method:
	// route.Put("/book", middleware.JWTProtected(), controllers.UpdateBook) // update one book by ID

	// Routes for DELETE method:
	// route.Delete("/book", middleware.JWTProtected(), controllers.DeleteBook) // delete one book by ID
}
