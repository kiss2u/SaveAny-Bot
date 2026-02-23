package web

import (
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/krau/SaveAny-Bot/storage"
)

// cacheEntry stores cached response with expiration time
type cacheEntry struct {
	data      interface{}
	expiration time.Time
}

// Simple in-memory cache for API responses
var (
	storageCache     cacheEntry
	storageCacheOnce sync.RWMutex
	cacheTTL         = 5 * time.Second // Cache TTL for 5 seconds
)

// getStorageCache returns cached storage data if still valid
func getStorageCache() ([]map[string]interface{}, bool) {
	storageCacheOnce.RLock()
	defer storageCacheOnce.RUnlock()
	
	if time.Now().Before(storageCache.expiration) {
		if data, ok := storageCache.data.([]map[string]interface{}); ok {
			return data, true
		}
	}
	return nil, false
}

// setStorageCache updates the storage cache
func setStorageCache(data []map[string]interface{}) {
	storageCacheOnce.Lock()
	defer storageCacheOnce.Unlock()
	
	storageCache.data = data
	storageCache.expiration = time.Now().Add(cacheTTL)
}

func (s *Server) handleGetStorages(c *fiber.Ctx) error {
	// Try to get from cache first
	if cached, ok := getStorageCache(); ok {
		return c.JSON(cached)
	}

	// Build storage list
	storages := make([]map[string]interface{}, 0)
	for name, st := range storage.Storages {
		storages = append(storages, map[string]interface{}{
			"name": name,
			"type": st.Type().String(),
		})
	}

	// Cache the result
	setStorageCache(storages)

	return c.JSON(storages)
}

type AddStorageRequest struct {
	Name     string                 `json:"name"`
	Type     string                 `json:"type"`
	Enable   bool                   `json:"enable"`
	Config   map[string]interface{} `json:"config"`
}

// invalidateStorageCache clears the storage cache
func invalidateStorageCache() {
	storageCacheOnce.Lock()
	defer storageCacheOnce.Unlock()
	storageCache = cacheEntry{}
}

func (s *Server) handleAddStorage(c *fiber.Ctx) error {
	var req AddStorageRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}

	// Invalidate cache when storage is added
	invalidateStorageCache()

	// TODO: Implement storage addition
	// This requires changes to config and storage packages

	return c.JSON(fiber.Map{"status": "ok", "message": "storage added"})
}

func (s *Server) handleDeleteStorage(c *fiber.Ctx) error {
	name := c.Params("name")
	if name == "" {
		return c.Status(400).JSON(fiber.Map{"error": "name required"})
	}

	// Invalidate cache when storage is deleted
	invalidateStorageCache()

	// TODO: Implement storage deletion

	return c.JSON(fiber.Map{"status": "ok", "message": "storage deleted"})
}
