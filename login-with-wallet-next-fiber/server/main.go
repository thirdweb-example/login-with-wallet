package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/thirdweb-dev/go-sdk/thirdweb"
)

type WalletLoginPayload struct {
	Payload struct {
		Payload struct {
			Domain         string    `json:"domain"`
			Address        string    `json:"address"`
			Nonce          string    `json:"nonce"`
			ExpirationTime time.Time `json:"expirationTime"`
			ChainId        int       `json:"chainId"`
		}
		Signature string `json:"signature"`
	} `json:"payload"`
}

func (w *WalletLoginPayload) ToThirdWeb() (*thirdweb.WalletLoginPayload, error) {
	signature, err := hex.DecodeString(w.Payload.Signature[2:])
	if err != nil {
		return nil, err
	}

	return &thirdweb.WalletLoginPayload{
		Payload: &thirdweb.WalletLoginPayloadData{
			Domain:         w.Payload.Payload.Domain,
			Address:        w.Payload.Payload.Address,
			Nonce:          w.Payload.Payload.Nonce,
			ExpirationTime: w.Payload.Payload.ExpirationTime,
			ChainId:        w.Payload.Payload.ChainId,
		},
		Signature: signature,
	}, nil
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	app := fiber.New()

	app.Post("/login", func(c *fiber.Ctx) error {
		w := new(WalletLoginPayload)

		if err := c.BodyParser(w); err != nil {
			return err
		}

		payload, err := w.ToThirdWeb()
		if err != nil {
			return err
		}

		sdk, err := thirdweb.NewThirdwebSDK("mumbai", &thirdweb.SDKOptions{
			PrivateKey: os.Getenv("ADMIN_PRIVATE_KEY"),
		})
		if err != nil {
			return err
		}

		domain := "thirdweb.com"
		token, err := sdk.Auth.GenerateAuthToken(domain, payload, nil)
		if err != nil {
			fmt.Println(err)
			return err
		}

		fmt.Println(token)

		return c.Status(200).JSON(&fiber.Map{
			"address": token,
		})
	})

	log.Fatal(app.Listen(":8000"))
}
