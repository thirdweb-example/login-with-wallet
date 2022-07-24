package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/thirdweb-dev/go-sdk/thirdweb"
)

type LoginPayload struct {
	Payload *thirdweb.WalletLoginPayload `json:"payload"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	app := fiber.New()

	app.Post("/login", func(c *fiber.Ctx) error {
		w := new(LoginPayload)

		if err := c.BodyParser(w); err != nil {
			return err
		}

		sdk, err := thirdweb.NewThirdwebSDK("mumbai", &thirdweb.SDKOptions{
			PrivateKey: os.Getenv("ADMIN_PRIVATE_KEY"),
		})
		if err != nil {
			return err
		}

		payload := w.Payload
		domain := "thirdweb.com"
		token, err := sdk.Auth.GenerateAuthToken(domain, payload, nil)
		if err != nil {
			return err
		}

		return c.Status(200).JSON(&fiber.Map{
			"address": token,
		})
	})

	log.Fatal(app.Listen(":8000"))
}
