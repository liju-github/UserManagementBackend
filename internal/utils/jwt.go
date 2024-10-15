package utils

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func GenerateJWT(email string, userID string, role string,expiry uint) (string, error) {
	claims := jwt.MapClaims{
		"email":  email,
		"userID": userID,
		"exp":    time.Now().Add(time.Hour * time.Duration(expiry)).Unix(),
		"role":   role,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte("secret"))
}


func JWTMiddleware(requiredRole string) fiber.Handler {
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

		// Check if the role matches the required role
		clientRole, ok := claims["role"].(string)
		if !ok {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid role in token"})
		}

		// Allow admin to access PUT, POST, DELETE methods without restrictions
		if clientRole == "admin" && (ctx.Method() == "PUT" || ctx.Method() == "POST" || ctx.Method() == "DELETE") {
			// Admin exception: allow PUT, POST, DELETE requests
			ctx.Locals("userID", claims["userID"])
			ctx.Locals("email", claims["email"])
			ctx.Locals("expiry", claims["exp"])
			ctx.Locals("role", claims["role"])
			log.Println("Admin access granted for request ", claims)
			return ctx.Next() // Admin can proceed
		}

		// Set user details in context locals for regular users
		ctx.Locals("userID", claims["userID"])
		ctx.Locals("email", claims["email"])
		ctx.Locals("expiry", claims["exp"])
		ctx.Locals("role", claims["role"])

		// Log the request claims
		log.Println("Request from ", claims)

		// Proceed to the next middleware or route handler
		return ctx.Next()
	}
}
