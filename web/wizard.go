package web

import (
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
)

type WizardStep int

const (
	WizardStepStart WizardStep = iota
	WizardStepBotToken
	WizardStepStorage
	WizardStepStorageConfig
	WizardStepUserID
	WizardStepConfirm
	WizardStepComplete
)

type WizardSession struct {
	Step       WizardStep      `json:"step"`
	Token      string          `json:"token,omitempty"`
	AppID      int            `json:"app_id,omitempty"`
	AppHash    string          `json:"app_hash,omitempty"`
	Storage    string          `json:"storage,omitempty"`
	StorageCfg map[string]interface{} `json:"storage_config,omitempty"`
	UserID     int64           `json:"user_id,omitempty"`
	CreatedAt  time.Time       `json:"created_at"`
}

var (
	wizardSessions = &sync.Map{}
	sessionTimeout = 30 * time.Minute
)

type WizardStartRequest struct {
	SessionID string `json:"session_id"`
}

func (s *Server) handleWizardStart(c *fiber.Ctx) error {
	var req WizardStartRequest
	if err := c.BodyParser(&req); err != nil {
		// Generate new session if not provided
		req.SessionID = generateSessionID()
	}

	session := &WizardSession{
		Step:      WizardStepBotToken,
		CreatedAt: time.Now(),
	}
	wizardSessions.Store(req.SessionID, session)

	return c.JSON(fiber.Map{
		"session_id": req.SessionID,
		"step":       int(session.Step),
		"message":    "Enter your Bot Token from @BotFather",
	})
}

func (s *Server) handleWizardStep(c *fiber.Ctx) error {
	sessionID := c.Query("session_id")
	if sessionID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "session_id required"})
	}

	sessionI, ok := wizardSessions.Load(sessionID)
	if !ok {
		return c.Status(404).JSON(fiber.Map{"error": "session not found"})
	}
	session := sessionI.(*WizardSession)

	// Check timeout
	if time.Since(session.CreatedAt) > sessionTimeout {
		wizardSessions.Delete(sessionID)
		return c.Status(410).JSON(fiber.Map{"error": "session expired"})
	}

	return c.JSON(fiber.Map{
		"step":    int(session.Step),
		"session": session,
	})
}

type WizardNextRequest struct {
	SessionID string                 `json:"session_id"`
	Data      map[string]interface{} `json:"data"`
}

func (s *Server) handleWizardNext(c *fiber.Ctx) error {
	var req WizardNextRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}

	sessionI, ok := wizardSessions.Load(req.SessionID)
	if !ok {
		return c.Status(404).JSON(fiber.Map{"error": "session not found"})
	}
	session := sessionI.(*WizardSession)

	// Process data based on current step
	switch session.Step {
	case WizardStepBotToken:
		if token, ok := req.Data["token"].(string); ok && token != "" {
			session.Token = token
			session.Step = WizardStepStorage
		} else {
			return c.Status(400).JSON(fiber.Map{"error": "token required"})
		}

	case WizardStepStorage:
		if storage, ok := req.Data["storage"].(string); ok && storage != "" {
			session.Storage = storage
			session.Step = WizardStepStorageConfig
		} else {
			return c.Status(400).JSON(fiber.Map{"error": "storage required"})
		}

	case WizardStepStorageConfig:
		if cfg, ok := req.Data["config"].(map[string]interface{}); ok {
			session.StorageCfg = cfg
			session.Step = WizardStepUserID
		} else {
			return c.Status(400).JSON(fiber.Map{"error": "config required"})
		}

	case WizardStepUserID:
		if userID, ok := req.Data["user_id"].(float64); ok {
			session.UserID = int64(userID)
			session.Step = WizardStepConfirm
		} else {
			return c.Status(400).JSON(fiber.Map{"error": "user_id required"})
		}

	case WizardStepConfirm:
		// Generate config file
		session.Step = WizardStepComplete
		// TODO: Generate and save config.toml

	default:
		return c.Status(400).JSON(fiber.Map{"error": "invalid step"})
	}

	wizardSessions.Store(req.SessionID, session)

	return c.JSON(fiber.Map{
		"step":    int(session.Step),
		"message": getStepMessage(session.Step),
	})
}

func getStepMessage(step WizardStep) string {
	switch step {
	case WizardStepBotToken:
		return "Enter your Bot Token from @BotFather"
	case WizardStepStorage:
		return "Select a storage backend"
	case WizardStepStorageConfig:
		return "Configure the storage"
	case WizardStepUserID:
		return "Enter your Telegram User ID"
	case WizardStepConfirm:
		return "Review and confirm"
	case WizardStepComplete:
		return "Setup complete! Bot is ready to use."
	default:
		return ""
	}
}

func generateSessionID() string {
	// Simple session ID generation
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
	}
	return string(b)
}
