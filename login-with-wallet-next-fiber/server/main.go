package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

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

		privateKey := os.Getenv("ADMIN_PRIVATE_KEY")
		if privateKey == "" {
			fmt.Printf("Missing ADMIN_PRIVATE_KEY environment variable")
			return errors.New("Admin private key not set")
		}

		sdk, err := thirdweb.NewThirdwebSDK("mumbai", &thirdweb.SDKOptions{
			PrivateKey: privateKey,
		})
		if err != nil {
			return err
		}

		// Get signed login payload from the frontend
		payload := w.Payload

		// Generate an access token with the SDK using the signed payload
		domain := "thirdweb.com"
		token, err := sdk.Auth.GenerateAuthToken(domain, payload, nil)
		if err != nil {
			return err
		}

		// Securely set httpOnly cookie on request to prevent XSS on frontend
		// And set path to / to enable access_token usage on all endpoints
		c.Cookie(&fiber.Cookie{
			Name:     "access_token",
			Value:    token,
			Path:     "/",
			HTTPOnly: true,
			Secure:   true,
			SameSite: "strict",
		})

		return c.SendStatus(200)
	})

	app.Post("/authenticate", func(c *fiber.Ctx) error {
		privateKey := os.Getenv("ADMIN_PRIVATE_KEY")
		if privateKey == "" {
			fmt.Printf("Missing ADMIN_PRIVATE_KEY environment variable")
			return errors.New("Admin private key not set")
		}

		sdk, err := thirdweb.NewThirdwebSDK("mumbai", &thirdweb.SDKOptions{
			PrivateKey: privateKey,
		})
		if err != nil {
			return err
		}

		// Get access token off cookies
		token := c.Cookies("access_token")
		if token == "" {
			return c.SendStatus(401)
		}

		// Authenticate token with the SDK
		domain := "thirdweb.com"
		address, err := sdk.Auth.Authenticate(domain, token)
		if err != nil {
			return c.SendStatus(401)
		}

		return c.Status(200).JSON(address)
	})

	app.Post("/logout", func(c *fiber.Ctx) error {
		c.Cookie(&fiber.Cookie{
			Name:    "access_token",
			Value:   "none",
			Expires: time.Unix(time.Now().Unix()+5*1000, 0),
		})
		return c.SendStatus(200)
	})

	log.Fatal(app.Listen(":8000"))
}
