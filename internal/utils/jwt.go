package utils

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/liju-github/user-management/internal/repository"
)

func GenerateJWT(email string, ID string, role string, expiry uint) (string, error) {
	claims := jwt.MapClaims{
		"email": email,
		"ID":    ID,
		"exp":   time.Now().Add(time.Hour * time.Duration(expiry)).Unix(),
		"role":  role,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte("secret"))
}

func JWTMiddleware(requiredRole string, userDB *repository.UserRepository) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		authHeader := ctx.Get("Authorization")
		if authHeader == "" {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing or invalid token"})
		}

		tokenStr := authHeader[len("Bearer "):]
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		})

		// Check if the token is valid and not expired
		if err != nil || !token.Valid {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
		}

		// Extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || claims["exp"] == nil {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token claims"})
		}

		// Compare the expiration time
		exp := claims["exp"].(float64) // JWT expiration is a float64
		if time.Now().Unix() > int64(exp) {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token has expired"})
		}

		// Extract role from JWT claims
		clientRole, ok := claims["role"].(string)
		if !ok {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid role in token"})
		}

		// Allow admin to access all routes except GET
		if clientRole == "admin" {
			ctx.Locals("ID", claims["ID"])
			ctx.Locals("email", claims["email"])
			ctx.Locals("expiry", claims["exp"])
			ctx.Locals("role", claims["role"])
			log.Println("Admin access granted for request ", claims)
			return ctx.Next() // Admin can proceed
		}

		// If a requiredRole is provided, ensure the user has that role
		if requiredRole != "" && requiredRole != clientRole {
			return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Insufficient privileges"})
		}


		userInfo,err := userDB.FindUserByID(claims["ID"].(string))
		if err!=nil {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Please try again"})
		}
		if userInfo.IsBlocked{
			return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "user is blocked"})
		}

		// Set user details in context locals
		ctx.Locals("ID", claims["ID"])
		ctx.Locals("email", claims["email"])
		ctx.Locals("expiry", claims["exp"])
		ctx.Locals("role", claims["role"])

		// Log the request claims
		log.Println("Request from ", claims)

		// Proceed to the next middleware or route handler
		return ctx.Next()
	}
}
