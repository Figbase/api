package main

import (
	"github.com/MTNNigeria/odyssey/pkg/middleware"
	"github.com/MTNNigeria/odyssey/pkg/routes"
	"github.com/MTNNigeria/odyssey/platform/database"
	"github.com/gofiber/fiber/v2"
)

// @title Odyssey API
// @version 1.0
// @description This is an auto-generated API Docs.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email zyro.kamson@mtn.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath /api
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	// Database connections.
	database.ConnectDb()

	// Define a new Fiber app.
	app := fiber.New()

	// Middlewares.
	middleware.FiberMiddleware(app) // Register Fiber's middleware for app.

	// Routes.
	routes.PublicRoutes(app)  // Register public routes for app.
	routes.PrivateRoutes(app) // Register private routes for app.
	routes.NotFoundRoute(app) // Register route for 404 Error.

	// Start server (with or without graceful shutdown).
	app.Listen(":3000")
}
