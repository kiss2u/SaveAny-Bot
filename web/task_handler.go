package web

import (
	"github.com/gofiber/fiber/v2"
	"github.com/krau/SaveAny-Bot/core"
)

func (s *Server) handleGetTasks(c *fiber.Ctx) error {
	running := core.GetRunningTasks(s.ctx)
	queued := core.GetQueuedTasks(s.ctx)

	type TaskInfo struct {
		ID        string `json:"id"`
		Title     string `json:"title"`
		Status    string `json:"status"`
		Created   int64  `json:"created"` // Unix timestamp
		Cancelled bool   `json:"cancelled"`
	}

	runningTasks := make([]TaskInfo, 0, len(running))
	for _, t := range running {
		runningTasks = append(runningTasks, TaskInfo{
			ID:        t.ID,
			Title:     t.Title,
			Status:    "running",
			Created:   t.Created.Unix(),
			Cancelled: t.Cancelled,
		})
	}

	queuedTasks := make([]TaskInfo, 0, len(queued))
	for _, t := range queued {
		queuedTasks = append(queuedTasks, TaskInfo{
			ID:        t.ID,
			Title:     t.Title,
			Status:    "queued",
			Created:   t.Created.Unix(),
			Cancelled: t.Cancelled,
		})
	}

	return c.JSON(fiber.Map{
		"running": runningTasks,
		"queued":  queuedTasks,
	})
}

type CancelTaskRequest struct {
	ID string `json:"id"`
}

func (s *Server) handleCancelTask(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{"error": "id required"})
	}

	err := core.CancelTask(s.ctx, id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"status": "ok"})
}
