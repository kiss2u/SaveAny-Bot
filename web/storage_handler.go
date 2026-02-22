package web

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kiss2u/SaveAny-Bot/storage"
)

func (s *Server) handleGetStorages(c *fiber.Ctx) error {
	storages := make([]map[string]interface{}, 0)
	for name, st := range storage.Storages {
		storages = append(storages, map[string]interface{}{
			"name": name,
			"type": st.Type().String(),
		})
	}
	return c.JSON(storages)
}

type AddStorageRequest struct {
	Name     string                 `json:"name"`
	Type     string                 `json:"type"`
	Enable   bool                   `json:"enable"`
	Config   map[string]interface{} `json:"config"`
}

func (s *Server) handleAddStorage(c *fiber.Ctx) error {
	var req AddStorageRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}

	// TODO: Implement storage addition
	// This requires changes to config and storage packages

	return c.JSON(fiber.Map{"status": "ok", "message": "storage added"})
}

func (s *Server) handleDeleteStorage(c *fiber.Ctx) error {
	name := c.Params("name")
	if name == "" {
		return c.Status(400).JSON(fiber.Map{"error": "name required"})
	}

	// TODO: Implement storage deletion

	return c.JSON(fiber.Map{"status": "ok", "message": "storage deleted"})
}
