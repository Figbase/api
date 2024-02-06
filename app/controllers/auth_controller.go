package controllers

import (
	"context"

	"time"

	"github.com/Figbase/api/app/models"
	"github.com/Figbase/api/pkg/utils"

	"github.com/Figbase/api/platform/cache"
	"github.com/Figbase/api/platform/database"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// UserSignUp method to create a new user.
// @Description Create a new user.
// @Summary create a new user
// @Tags User
// @Accept json
// @Produce json
// @Param email body string true "Email"
// @Param password body string true "Password"
// @Param user_role body string true "User role"
// @Success 200 {object} models.User
// @Router /v1/auth/signup [post]
func UserSignUp(c *fiber.Ctx) error {
	// Create a new user auth struct.
	signUp := &models.SignUp{}

	// Checking received data from JSON body.
	if err := c.BodyParser(signUp); err != nil {
		// Return status 400 and error message.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	// Check if the email has been used to signup before.
	existingUser := &models.User{}
	if err := database.DB.Db.Where("email = ?", signUp.Email).First(existingUser).Error; err == nil {
		// Return status 400 and error message indicating that the email is already registered.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Email address is already registered",
		})
	}

	// Create a new validator for a User model.
	validate := utils.NewValidator()

	// Validate sign up fields.
	if err := validate.Struct(signUp); err != nil {
		// Return, if some fields are not valid.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": utils.ValidatorErrors(err),
		})
	}

	// Checking role from sign up data.
	role, err := utils.VerifyRole("user")
	if err != nil {
		// Return status 400 and error message.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	// Create a new user struct.
	user := &models.User{}

	// Set initialized default data for user:
	user.ID = uuid.New()
	user.CreatedAt = time.Now()
	user.Email = signUp.Email
	user.PasswordHash = utils.GeneratePassword(signUp.Password)
	user.UserStatus = 1 // 0 == blocked, 1 == active
	user.UserRole = role

	// Validate user fields.
	if err := validate.Struct(user); err != nil {
		// Return, if some fields are not valid.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": utils.ValidatorErrors(err),
		})
	}

	// Checking received data from JSON body.
	if err := c.BodyParser(user); err != nil {
		// Return status 400 and error message.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	// Create a new user with validated data.
	database.DB.Db.Create(&user)

	// Get role credentials from created user.
	credentials, err := utils.GetCredentialsByRole(user.UserRole)
	if err != nil {
		// Return status 400 and error message.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	// Generate a new pair of access and refresh tokens.
	tokens, err := utils.GenerateNewTokens(user.ID.String(), credentials)
	if err != nil {
		// Return status 500 and token generation error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	// Create a new Redis connection.
	connRedis, err := cache.RedisConnection()
	if err != nil {
		// Return status 500 and Redis connection error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	// Save refresh token to Redis.
	errSaveToRedis := connRedis.Set(context.Background(), user.ID.String(), tokens.Refresh, 0).Err()
	if errSaveToRedis != nil {
		// Return status 500 and Redis connection error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": errSaveToRedis.Error(),
		})
	}

	// Delete password hash field from JSON view.
	user.PasswordHash = ""

	// Return status 200 OK with user details and tokens.
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": nil,
		"user":    user,
		"tokens": fiber.Map{
			"access":  tokens.Access,
			"refresh": tokens.Refresh,
		},
	})
}

// UserSignIn method to auth user and return access and refresh tokens.
// @Description Auth user and return access and refresh token.
// @Summary auth user and return access and refresh token
// @Tags User
// @Accept json
// @Produce json
// @Param email body string true "User Email"
// @Param password body string true "User Password"
// @Success 200 {string} status "ok"
// @Router /v1/auth/signin [post]
func UserSignIn(c *fiber.Ctx) error {
	// Create a new user auth struct.
	signIn := &models.SignIn{}

	// Checking received data from JSON body.
	if err := c.BodyParser(signIn); err != nil {
		// Return status 400 and error message.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	// Create a new user struct.
	user := &models.User{}

	// Get user by email.
	result := database.DB.Db.Where("email = ?", signIn.Email).First(&user)

	// Check if the user was not found.
	if result.Error != nil {
		// Return, if user not found.
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "User with the given email is not found",
		})
	}

	// Compare given user password with stored in found user.
	compareUserPassword := utils.ComparePasswords(user.PasswordHash, signIn.Password)
	if !compareUserPassword {
		// Return, if password is not compare to stored in database.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Wrong user email address or password",
		})
	}

	// Get role credentials from founded user.
	credentials, err := utils.GetCredentialsByRole(user.UserRole)
	if err != nil {
		// Return status 400 and error message.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	// Generate a new pair of access and refresh tokens.
	tokens, err := utils.GenerateNewTokens(user.ID.String(), credentials)
	if err != nil {
		// Return status 500 and token generation error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	// Define user ID.
	userID := user.ID.String()

	// Create a new Redis connection.
	connRedis, err := cache.RedisConnection()
	if err != nil {
		// Return status 500 and Redis connection error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	// Save refresh token to Redis.
	errSaveToRedis := connRedis.Set(context.Background(), userID, tokens.Refresh, 0).Err()
	if errSaveToRedis != nil {
		// Return status 500 and Redis connection error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": errSaveToRedis.Error(),
		})
	}

	// Delete password hash field from JSON view.
	user.PasswordHash = ""

	// Return status 200 OK.
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": nil,
		"user":    user,
		"tokens": fiber.Map{
			"access":  tokens.Access,
			"refresh": tokens.Refresh,
		},
	})
}

// UserSignOut method to de-authorize user and delete refresh token from Redis.
// @Description De-authorize user and delete refresh token from Redis.
// @Summary de-authorize user and delete refresh token from Redis
// @Tags User
// @Accept json
// @Produce json
// @Success 204 {string} status "ok"
// @Security ApiKeyAuth
// @Router /v1/auth/signout [post]
func UserSignOut(c *fiber.Ctx) error {
	// Get claims from JWT.
	claims, err := utils.ExtractTokenMetadata(c)
	if err != nil {
		// Return status 500 and JWT parse error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	// Define user ID.
	userID := claims.UserID.String()

	// Create a new Redis connection.
	connRedis, err := cache.RedisConnection()
	if err != nil {
		// Return status 500 and Redis connection error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	// Save refresh token to Redis.
	errDelFromRedis := connRedis.Del(context.Background(), userID).Err()
	if errDelFromRedis != nil {
		// Return status 500 and Redis deletion error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": errDelFromRedis.Error(),
		})
	}

	// Return status 204 no content.
	return c.SendStatus(fiber.StatusNoContent)
}
