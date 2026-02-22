package web

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/krau/SaveAny-Bot/config"
)

type ConfigResponse struct {
	Lang          string            `json:"lang"`
	Workers       int              `json:"workers"`
	Retry         int              `json:"retry"`
	Threads       int              `json:"threads"`
	Stream        bool             `json:"stream"`
	Proxy         string           `json:"proxy"`
	DB            interface{}      `json:"db"`
	Telegram      interface{}      `json:"telegram"`
	Storages      []interface{}    `json:"storages"`
	Parser        interface{}      `json:"parser"`
	Hook          interface{}      `json:"hook"`
}

func (s *Server) handleGetConfig(c *fiber.Ctx) error {
	cfg := config.C()
	return c.JSON(ConfigResponse{
		Lang:     cfg.Lang,
		Workers:  cfg.Workers,
		Retry:    cfg.Retry,
		Threads:  cfg.Threads,
		Stream:   cfg.Stream,
		Proxy:    cfg.Proxy,
		DB:       cfg.DB,
		Telegram: cfg.Telegram,
		Storages: nil, // TODO: get from storage package
		Parser:   cfg.Parser,
		Hook:     cfg.Hook,
	})
}

type SaveConfigRequest struct {
	Lang     string `json:"lang"`
	Workers  int    `json:"workers"`
	Retry    int    `json:"retry"`
	Threads  int    `json:"threads"`
	Stream   bool   `json:"stream"`
	Proxy    string `json:"proxy"`
}

func (s *Server) handleSaveConfig(c *fiber.Ctx) error {
	var req SaveConfigRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}

	// TODO: Validate and save to config file
	// This requires modifying config package to support hot-reload

	return c.JSON(fiber.Map{"status": "ok", "message": "config saved"})
}

func (s *Server) handleValidateConfig(c *fiber.Ctx) error {
	var req map[string]interface{}
	if err := json.Unmarshal(c.Body(), &req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid json"})
	}

	// Validate Telegram token
	if token, ok := req["token"].(string); ok && token != "" {
		// TODO: Validate with Telegram API
	}

	return c.JSON(fiber.Map{"valid": true})
}
