package notify

import (
	"sync"

	"github.com/celestix/gotgproto/ext"
)

type AdminNotifier struct {
	ctx      *ext.Context
	adminIDs []int64
	mu       sync.RWMutex
	enabled  bool
}

var notifier *AdminNotifier

func NewAdminNotifier(ctx *ext.Context, adminIDs []int64) *AdminNotifier {
	return &AdminNotifier{
		ctx:      ctx,
		adminIDs: adminIDs,
		enabled:  len(adminIDs) > 0 && ctx != nil,
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
	if n.ctx == nil {
		return
	}

	// Use ext.ReplyTextString to send a simple text message
	n.ctx.SendMessage(chatID, &ext.TextMessage{
		Text: msg,
	})
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
