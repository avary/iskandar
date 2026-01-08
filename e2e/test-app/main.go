package main

import (
	"bufio"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

func main() {
	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Static("/assets", "./public")

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{
			"Title": "Hello, World!",
		})
	})

	app.Post("/api/data", func(c *fiber.Ctx) error {
		data := struct {
			Name    string `json:"name"`
			Surname string `json:"surname"`
		}{
			Name:    "Peter",
			Surname: "Perlepes",
		}
		return c.JSON(data)
	})
	// Query parameters test
	app.Get("/search", func(c *fiber.Ctx) error {
		query := c.Query("q", "")
		limit := c.Query("limit", "10")
		page := c.Query("page", "1")

		return c.JSON(fiber.Map{
			"query":   query,
			"limit":   limit,
			"page":    page,
			"results": []string{"result1", "result2", "result3"},
		})
	})

	// 404 Error
	app.Get("/not-found", func(c *fiber.Ctx) error {
		return c.Status(404).JSON(fiber.Map{
			"error":   "Not Found",
			"message": "The requested resource does not exist",
		})
	})

	// 500 Internal Server Error
	app.Get("/server-error", func(c *fiber.Ctx) error {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Internal Server Error",
			"message": "Something went wrong on the server",
		})
	})

	// 400 Bad Request
	app.Get("/bad-request", func(c *fiber.Ctx) error {
		return c.Status(400).JSON(fiber.Map{
			"error":   "Bad Request",
			"message": "Invalid request parameters",
		})
	})

	// Echo all request headers
	app.Get("/headers/echo", func(c *fiber.Ctx) error {
		headers := make(map[string]string)
		c.Request().Header.VisitAll(func(key, value []byte) {
			headers[string(key)] = string(value)
		})
		return c.JSON(fiber.Map{
			"headers": headers,
		})
	})

	// Test Authorization header
	app.Get("/headers/auth", func(c *fiber.Ctx) error {
		auth := c.Get("Authorization")
		if auth == "" {
			return c.Status(401).JSON(fiber.Map{
				"error":   "Unauthorized",
				"message": "Authorization header missing",
			})
		}
		return c.JSON(fiber.Map{
			"authorized":  true,
			"auth_header": auth,
		})
	})

	// Test custom response headers
	app.Get("/headers/response", func(c *fiber.Ctx) error {
		c.Set("X-Custom-Header", "CustomValue")
		c.Set("X-Server", "TestApp")
		c.Set("Cache-Control", "no-cache")
		return c.JSON(fiber.Map{
			"message": "Check response headers",
		})
	})

	// Streaming response with multiple writes
	app.Get("/stream", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/plain")
		c.Set("Transfer-Encoding", "chunked")
		c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
			for i := 1; i <= 5; i++ {
				fmt.Fprintf(w, "Chunk %d: This is streaming data\n", i)
				w.Flush()
				time.Sleep(1 * time.Second)
			}
			fmt.Fprintf(w, "Stream complete!\n")
		})
		return nil
	})

	app.Listen(":3003")
}
