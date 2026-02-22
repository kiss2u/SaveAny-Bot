package handlers

import (
	"context"
	"fmt"
	"strings"

	"github.com/celestix/gotgproto/ext"
	"github.com/celestix/gotgproto/ext/utils"
	"github.com/celestix/gotgproto/tg"
	"github.com/krau/SaveAny-Bot/config"
	"github.com/krau/SaveAny-Bot/core"
	"github.com/krau/SaveAny-Bot/storage"
)

// MenuCallback prefixes
const (
	MenuCallbackStatus    = "menu:status"
	MenuCallbackTasks    = "menu:tasks"
	MenuCallbackStorages = "menu:storages"
	MenuCallbackSilent   = "menu:silent"
	MenuCallbackRefresh  = "menu:refresh"
)

func handleMenuCmd(ctx *ext.Context, u *ext.Update) error {
	return showMainMenu(ctx, u.EffectiveChat.ID)
}

func showMainMenu(ctx *ext.Context, chatID int64, msgID ...int) error {
	// Get running tasks count
	runningTasks := core.GetRunningTasks(context.Background())
	queuedTasks := core.GetQueuedTasks(context.Background())

	statusText := fmt.Sprintf("📊 *SaveAny-Bot Status*\n\n✅ Bot: Online\n📥 Running: %d\n⏳ Queued: %d\n💾 Storages: %d\n\n_Use buttons below or send commands_",
		len(runningTasks),
		len(queuedTasks),
		len(storage.Storages),
	)

	// Build inline keyboard
	rows := [][]utils.InlineKeyboardButton{
		{
			{Text: "📊 Status", CallbackData: MenuCallbackStatus},
			{Text: "📋 Tasks", CallbackData: MenuCallbackTasks},
		},
		{
			{Text: "💾 Storages", CallbackData: MenuCallbackStorages},
			{Text: "🔇 Silent Mode", CallbackData: MenuCallbackSilent},
		},
		{
			{Text: "🔄 Refresh", CallbackData: MenuCallbackRefresh},
		},
	}

	keyboard := utils.InlineKeyboardMarkup{
		InlineKeyboard: rows,
	}

	// Send menu message
	if len(msgID) > 0 {
		_, err := ctx.EditMessage(chatID, &tg.MessagesEditMessageRequest{
			Message:    statusText,
			MsgID:      msgID[0],
			ReplyMarkup: &keyboard,
		})
		return err
	}

	_, err := ctx.SendMessage(chatID, &tg.MessagesSendMessageRequest{
		Message:    statusText,
		ReplyMarkup: &keyboard,
	})
	return err
}

func handleMenuCallback(ctx *ext.Context, u *ext.Update) error {
	callbackData := string(u.CallbackQuery.Data)
	chatID := u.CallbackQuery.GetUserID()
	msgID := u.CallbackQuery.GetMsgID()

	// Answer callback first
	ctx.AnswerCallbackQuery(&tg.MessagesSetBotCallbackAnswerRequest{
		QueryID: u.CallbackQuery.QueryID,
	})

	switch {
	case strings.HasPrefix(callbackData, MenuCallbackStatus):
		return showStatusCallback(ctx, chatID, msgID)
	case strings.HasPrefix(callbackData, MenuCallbackTasks):
		return showTasksCallback(ctx, chatID, msgID)
	case strings.HasPrefix(callbackData, MenuCallbackStorages):
		return showStoragesCallback(ctx, chatID, msgID)
	case strings.HasPrefix(callbackData, MenuCallbackSilent):
		return toggleSilentCallback(ctx, chatID, msgID)
	case strings.HasPrefix(callbackData, MenuCallbackRefresh):
		return showMainMenu(ctx, chatID, msgID)
	}

	return nil
}

func showStatusCallback(ctx *ext.Context, chatID int64, msgID int) error {
	runningTasks := core.GetRunningTasks(context.Background())
	queuedTasks := core.GetQueuedTasks(context.Background())

	// Get bot version
	shortHash := config.GitCommit
	if len(shortHash) > 7 {
		shortHash = shortHash[:7]
	}

	statusText := fmt.Sprintf(`📊 *System Status*

✅ **Bot Status**: Running
📥 **Running Tasks**: %d
⏳ **Queued Tasks**: %d
💾 **Active Storages**: %d
⚙️ **Workers**: %d
🔄 **Version**: %s (%s)

_Updated: just now_`,
		len(runningTasks),
		len(queuedTasks),
		len(storage.Storages),
		config.C().Workers,
		config.Version,
		shortHash,
	)

	// Add back button
	rows := [][]utils.InlineKeyboardButton{
		{
			{Text: "🔙 Back to Menu", CallbackData: MenuCallbackRefresh},
		},
	}

	keyboard := utils.InlineKeyboardMarkup{
		InlineKeyboard: rows,
	}

	_, err := ctx.EditMessage(chatID, &tg.MessagesEditMessageRequest{
		Message:    statusText,
		MsgID:      msgID,
		ReplyMarkup: &keyboard,
	})
	return err
}

func showTasksCallback(ctx *ext.Context, chatID int64, msgID int) error {
	runningTasks := core.GetRunningTasks(context.Background())
	queuedTasks := core.GetQueuedTasks(context.Background())

	var tasksText string

	if len(runningTasks) == 0 && len(queuedTasks) == 0 {
		tasksText = "📋 *Tasks*\n\nNo active tasks"
	} else {
		tasksText = "📋 *Tasks*\n\n"

		if len(runningTasks) > 0 {
			tasksText += "📥 *Running:*\n"
			for i, task := range runningTasks {
				if i >= 5 { // Show max 5 tasks
					tasksText += fmt.Sprintf("\n... and %d more", len(runningTasks)-5)
					break
				}
				tasksText += fmt.Sprintf("• %s\n", task.Title)
			}
		}

		if len(queuedTasks) > 0 {
			tasksText += "\n⏳ *Queued:*\n"
			for i, task := range queuedTasks {
				if i >= 5 {
					tasksText += fmt.Sprintf("\n... and %d more", len(queuedTasks)-5)
					break
				}
				tasksText += fmt.Sprintf("• %s\n", task.Title)
			}
		}
	}

	// Add action buttons
	rows := [][]utils.InlineKeyboardButton{
		{
			{Text: "🔄 Refresh", CallbackData: MenuCallbackTasks},
			{Text: "🔙 Back", CallbackData: MenuCallbackRefresh},
		},
	}

	keyboard := utils.InlineKeyboardMarkup{
		InlineKeyboard: rows,
	}

	_, err := ctx.EditMessage(chatID, &tg.MessagesEditMessageRequest{
		Message:    tasksText,
		MsgID:      msgID,
		ReplyMarkup: &keyboard,
	})
	return err
}

func showStoragesCallback(ctx *ext.Context, chatID int64, msgID int) error {
	storagesText := "💾 *Storages*\n\n"

	for name, s := range storage.Storages {
		storType := s.Type().String()
		storagesText += fmt.Sprintf("• *%s* (%s)\n", name, storType)
	}

	if len(storage.Storages) == 0 {
		storagesText += "_No storages configured_"
	}

	// Add back button
	rows := [][]utils.InlineKeyboardButton{
		{
			{Text: "🔙 Back to Menu", CallbackData: MenuCallbackRefresh},
		},
	}

	keyboard := utils.InlineKeyboardMarkup{
		InlineKeyboard: rows,
	}

	_, err := ctx.EditMessage(chatID, &tg.MessagesEditMessageRequest{
		Message:    storagesText,
		MsgID:      msgID,
		ReplyMarkup: &keyboard,
	})
	return err
}

func toggleSilentCallback(ctx *ext.Context, chatID int64, msgID int) error {
	// Get user's silent mode status from database or config
	userID := chatID
	// For now, just show current status - actual toggle would need database access

	silentText := "🔇 *Silent Mode*\n\n"
	silentText += "Current: _Use /silent to toggle_\n\n"
	silentText += "When enabled, bot won't send completion notifications for downloads."

	// Add back button
	rows := [][]utils.InlineKeyboardButton{
		{
			{Text: "🔙 Back to Menu", CallbackData: MenuCallbackRefresh},
		},
	}

	keyboard := utils.InlineKeyboardMarkup{
		InlineKeyboard: rows,
	}

	_, err := ctx.EditMessage(chatID, &tg.MessagesEditMessageRequest{
		Message:    silentText,
		MsgID:      msgID,
		ReplyMarkup: &keyboard,
	})
	return err
}
