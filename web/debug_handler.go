package web

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/krau/SaveAny-Bot/database"
)

type MessageLogResponse struct {
	ID          uint   `json:"id"`
	ChatID      int64  `json:"chat_id"`
	UserID      int64  `json:"user_id"`
	Message     string `json:"message"`
	MessageType string `json:"message_type"`
	CreatedAt   string `json:"created_at"`
}

// handleGetMessageLogs returns recent message logs
func (s *Server) handleGetMessageLogs(c *fiber.Ctx) error {
	limitStr := c.Query("limit", "50")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit > 100 {
		limit = 50
	}

	logs, err := database.GetMessageLogs(s.ctx, limit)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	result := make([]MessageLogResponse, 0, len(logs))
	for _, log := range logs {
		result = append(result, MessageLogResponse{
			ID:          log.ID,
			ChatID:      log.ChatID,
			UserID:      log.UserID,
			Message:     log.Message,
			MessageType: log.MessageType,
			CreatedAt:   log.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return c.JSON(result)
}

// handleGetMessageStats returns message statistics
func (s *Server) handleGetMessageStats(c *fiber.Ctx) error {
	stats, err := database.GetMessageStats(s.ctx)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(stats)
}

// handleClearMessageLogs clears all message logs
func (s *Server) handleClearMessageLogs(c *fiber.Ctx) error {
	err := database.ClearMessageLogs(s.ctx)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"status": "ok", "message": "message logs cleared"})
}
