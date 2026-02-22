package web

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/charmbracelet/log"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/kiss2u/SaveAny-Bot/config"
)

type Server struct {
	app    *fiber.App
	config *config.WebConfig
	ctx    context.Context
	mu     sync.RWMutex
}

var serverInstance *Server

func New(ctx context.Context, cfg *config.WebConfig) *Server {
	if serverInstance != nil {
		return serverInstance
	}

	app := fiber.New(fiber.Config{
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		AppName:      "SaveAny-Bot Web",
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New())

	serverInstance = &Server{
		app:    app,
		config: cfg,
		ctx:    ctx,
	}

	return serverInstance
}

func (s *Server) setupRoutes() {
	// Apply auth middleware if configured
	authMiddleware := NewAuthMiddleware(&AuthConfig{
		Username: s.config.Username,
		Password: s.config.Password,
	})

	// API routes
	api := s.app.Group("/api", authMiddleware)

	// Status
	api.Get("/status", s.handleStatus)
	api.Get("/health", s.handleHealth)

	// Config
	api.Get("/config", s.handleGetConfig)
	api.Post("/config", s.handleSaveConfig)
	api.Post("/config/validate", s.handleValidateConfig)

	// Storage
	api.Get("/storages", s.handleGetStorages)
	api.Post("/storages", s.handleAddStorage)
	api.Delete("/storages/:name", s.handleDeleteStorage)

	// Tasks
	api.Get("/tasks", s.handleGetTasks)
	api.Delete("/tasks/:id", s.handleCancelTask)

	// Wizard - 公开路由（不需要认证）
	wizard := s.app.Group("/api/wizard")
	wizard.Get("/fields", s.handleWizardStorageFields)
	wizard.Get("/step", s.handleWizardStep)
	wizard.Post("/step", s.handleWizardNext)
	wizard.Post("/start", s.handleWizardStart)
	wizard.Post("/generate", s.handleWizardGenerate)
	wizard.Get("/download", s.handleWizardDownload)

	// Static files
	s.app.Get("/", s.handleIndex)
	s.app.Static("/static", "./web/static")
}

func (s *Server) Run() error {
	s.setupRoutes()
	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
	log.Info("Starting web server", "addr", addr)
	return s.app.Listen(addr)
}

func (s *Server) RunWithTLS(certFile, keyFile string) error {
	s.setupRoutes()
	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
	log.Info("Starting web server with TLS", "addr", addr)
	return s.app.ListenTLS(addr, certFile, keyFile)
}

func (s *Server) Shutdown() error {
	log.Info("Shutting down web server")
	return s.app.Shutdown()
}
