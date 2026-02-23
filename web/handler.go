package web

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kiss2u/SaveAny-Bot/core"
)

type StatusResponse struct {
	BotConnected  bool          `json:"bot_connected"`
	RunningTasks  int           `json:"running_tasks"`
	QueuedTasks   int           `json:"queued_tasks"`
	Uptime        string        `json:"uptime"`
	Version       string        `json:"version"`
	StorageCount  int           `json:"storage_count"`
}

func (s *Server) handleStatus(c *fiber.Ctx) error {
	// TODO: Get actual bot connection status
	running := core.GetRunningTasks(s.ctx)
	queued := core.GetQueuedTasks(s.ctx)

	return c.JSON(StatusResponse{
		BotConnected: true,
		RunningTasks: len(running),
		QueuedTasks:  len(queued),
		Uptime:       "0h", // TODO
		Version:      "0.1.0",
		StorageCount: 0,
	})
}

func (s *Server) handleHealth(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status": "ok",
	})
}
