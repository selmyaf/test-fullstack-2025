package main

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/net/context"
)

var ctx = context.Background()

type User struct {
	Realname string `json:"realname"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func factorial(n int) int64 {
	if n == 0 || n == 1 {
		return 1
	}
	res := int64(1)
	for i := 2; i <= n; i++ {
		res *= int64(i)
	}
	return res
}

func f(n int) int64 {
	return int64(math.Ceil(float64(factorial(n)) / math.Pow(2, float64(n))))
}

func main() {
	fmt.Println("f(5) =", f(5)) 

	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	app := fiber.New()

	app.Post("/login", func(c *fiber.Ctx) error {
		type LoginRequest struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		var req LoginRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
		}

		key := fmt.Sprintf("login_%s", req.Username)
		val, err := rdb.Get(ctx, key).Result()
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "User not found"})
		}
		var user User
		if err := json.Unmarshal([]byte(val), &user); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Invalid user data"})
		}
		h := sha1.New()
		h.Write([]byte(req.Password))
		hashed := hex.EncodeToString(h.Sum(nil))

		if hashed != user.Password {
			return c.Status(401).JSON(fiber.Map{"error": "Invalid password"})
		}
		return c.JSON(fiber.Map{
			"message":  "Login success",
			"realname": user.Realname,
			"email":    user.Email,
		})
	})

	log.Fatal(app.Listen(":3000"))
}
