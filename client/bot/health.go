package bot

import (
	"context"
	"math"
	"sync"
	"time"

	"github.com/celestix/gotgproto"
	"github.com/charmbracelet/log"
)

type ConnectionStatus int

const (
	ConnectionStatusUnknown ConnectionStatus = iota
	ConnectionStatusConnected
	ConnectionStatusConnecting
	ConnectionStatusReconnecting
	ConnectionStatusDisconnected
)

type ConnectionCallback func(oldStatus, newStatus ConnectionStatus)

type HealthChecker struct {
	client        *gotgproto.Client
	mu            sync.RWMutex
	status        ConnectionStatus
	lastSuccess   time.Time
	lastError     error
	failCount     int
	checkInterval time.Duration
	maxRetries   int
	ctx          context.Context
	cancel       context.CancelFunc
	// Callbacks for status changes
	onDisconnected func()
	onReconnected  func()
	onReconnectFailed func()
}

var healthChecker *HealthChecker

func NewHealthChecker(client *gotgproto.Client, checkInterval time.Duration, maxRetries int) *HealthChecker {
	ctx, cancel := context.WithCancel(context.Background())
	return &HealthChecker{
		client:       client,
		status:       ConnectionStatusUnknown,
		checkInterval: checkInterval,
		maxRetries:   maxRetries,
		ctx:          ctx,
		cancel:       cancel,
	}
}

func (h *HealthChecker) Start(ctx context.Context) {
	h.mu.Lock()
	h.status = ConnectionStatusConnecting
	h.mu.Unlock()

	ticker := time.NewTicker(h.checkInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-h.ctx.Done():
			return
		case <-ticker.C:
			h.checkConnection()
		}
	}
}

func (h *HealthChecker) checkConnection() {
	err := h.client.Ping(h.ctx)
	h.mu.Lock()
	defer h.mu.Unlock()

	oldStatus := h.status

	if err != nil {
		h.failCount++
		h.lastError = err
		log.Warn("Health check failed", "error", err, "fail_count", h.failCount)

		if h.failCount >= 3 && h.status != ConnectionStatusReconnecting {
			h.status = ConnectionStatusReconnecting
			if h.onDisconnected != nil {
				go h.onDisconnected()
			}
			go h.handleReconnect()
		}
	} else {
		h.failCount = 0
		h.lastSuccess = time.Now()
		if h.status == ConnectionStatusReconnecting {
			h.status = ConnectionStatusConnected
			log.Info("Connection restored")
			if h.onReconnected != nil {
				go h.onReconnected()
			}
		} else if h.status == ConnectionStatusConnecting {
			h.status = ConnectionStatusConnected
		}
	}

	// Notify if status changed to disconnected
	if oldStatus != ConnectionStatusDisconnected && h.status == ConnectionStatusDisconnected {
		if h.onReconnectFailed != nil {
			go h.onReconnectFailed()
		}
	}
}

func (h *HealthChecker) handleReconnect() {
	logger := log.FromContext(h.ctx)
	logger.Warn("Connection unstable, attempting to reconnect...")

	for attempt := 1; attempt <= h.maxRetries; attempt++ {
		select {
		case <-h.ctx.Done():
			return
		default:
		}

		// Exponential backoff: 1s, 2s, 4s, 8s, 16s, 32s, 60s (capped)
		sleepDuration := time.Duration(math.Min(float64(attempt*attempt), 60)) * time.Second
		logger.Infof("Reconnect attempt %d/%d in %v", attempt, h.maxRetries, sleepDuration)
		time.Sleep(sleepDuration)

		if err := h.client.Ping(h.ctx); err == nil {
			h.mu.Lock()
			h.status = ConnectionStatusConnected
			h.failCount = 0
			h.mu.Unlock()
			logger.Info("Reconnected successfully")
			return
		}

		logger.Warnf("Reconnect attempt %d failed: %v", attempt, h.lastError)
	}

	h.mu.Lock()
	h.status = ConnectionStatusDisconnected
	h.mu.Unlock()

	logger.Error("Max reconnection attempts reached, manual intervention required")
}

func (h *HealthChecker) GetStatus() (ConnectionStatus, time.Time, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.status, h.lastSuccess, h.lastError
}

func (h *HealthChecker) Stop() {
	h.cancel()
}
