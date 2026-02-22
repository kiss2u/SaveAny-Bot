package notify

import (
	"sync"

	"github.com/celestix/gotgproto"
	"github.com/charmbracelet/log"
)

type AdminNotifier struct {
	client    *gotgproto.Client
	adminIDs  []int64
	mu        sync.RWMutex
	enabled   bool
}

var notifier *AdminNotifier

func NewAdminNotifier(client *gotgproto.Client, adminIDs []int64) *AdminNotifier {
	return &AdminNotifier{
		client:   client,
		adminIDs: adminIDs,
		enabled:  len(adminIDs) > 0 && client != nil,
	}
}

func (n *AdminNotifier) Notify(msg string) {
	if !n.enabled {
		return
	}

	n.mu.RLock()
	adminIDs := make([]int64, len(n.adminIDs))
	copy(adminIDs, n.adminIDs)
	n.mu.RUnlock()

	for _, id := range adminIDs {
		go n.sendMessage(id, msg)
	}
}

func (n *AdminNotifier) sendMessage(chatID int64, msg string) {
	if n.client == nil {
		return
	}

	_, err := n.client.SendMessage(chatID, msg, nil)
	if err != nil {
		log.Error("Failed to send admin notification", "chat_id", chatID, "error", err)
	}
}

func (n *AdminNotifier) NotifyDisconnected() {
	n.Notify("⚠️ Bot 断开连接，正在尝试重连...")
}

func (n *AdminNotifier) NotifyReconnected() {
	n.Notify("✅ Bot 重连成功")
}

func (n *AdminNotifier) NotifyReconnectFailed() {
	n.Notify("❌ Bot 重连失败，需要手动检查")
}

func (n *AdminNotifier) NotifyTaskFailed(taskTitle, errorMsg string) {
	n.Notify("❌ 任务失败: " + taskTitle + "\n错误: " + errorMsg)
}

func (n *AdminNotifier) NotifyTaskSuccess(taskTitle string) {
	n.Notify("✅ 任务完成: " + taskTitle)
}

func (n *AdminNotifier) NotifyStartup() {
	n.Notify("🚀 SaveAny-Bot 已启动")
}

func (n *AdminNotifier) NotifyShutdown() {
	n.Notify("👋 SaveAny-Bot 已关闭")
}
