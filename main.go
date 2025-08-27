package main

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)

var ctx = context.Background()

type User struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func hitungFaktorial(n int) int64 {
	if n <= 1 {
		return 1
	}
	result := int64(1)
	for i := 2; i <= n; i++ {
		result *= int64(i)
	}
	return result
}

func specialFunction(n int) int64 {
	fak := hitungFaktorial(n)
	return int64(math.Ceil(float64(fak) / math.Pow(2, float64(n))))
}

func main() {
	fmt.Println("Hasil f(5) =", specialFunction(5))
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	app := fiber.New()
	app.Post("/login", func(c *fiber.Ctx) error {
		type LoginInput struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		var loginData LoginInput
		if err := c.BodyParser(&loginData); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Request invalid"})
		}

		redisKey := fmt.Sprintf("user_%s", loginData.Username)
		storedUserJSON, err := redisClient.Get(ctx, redisKey).Result()
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "User not found"})
		}

		var storedUser User
		if err := json.Unmarshal([]byte(storedUserJSON), &storedUser); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "User data corrupted"})
		}

		hash := sha1.New()
		hash.Write([]byte(loginData.Password))
		hashedPassword := hex.EncodeToString(hash.Sum(nil))

		if hashedPassword != storedUser.Password {
			return c.Status(401).JSON(fiber.Map{"error": "Password salah"})
		}

		return c.JSON(fiber.Map{
			"message": "Login berhasil",
			"name":    storedUser.Name,
			"email":   storedUser.Email,
		})
	})

	log.Fatal(app.Listen(":3000"))
}
