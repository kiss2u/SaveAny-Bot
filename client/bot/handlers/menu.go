package handlers

import (
	"context"
	"fmt"
	"strings"

	"github.com/celestix/gotgproto/ext"
	"github.com/gotd/td/tg"
	"github.com/krau/SaveAny-Bot/config"
	"github.com/krau/SaveAny-Bot/core"
	"github.com/krau/SaveAny-Bot/storage"
)

// MenuCallback prefixes
const (
	MenuCallbackStatus    = "menu:status"
	MenuCallbackTasks    = "menu:tasks"
	MenuCallbackStorages = "menu:storages"
	MenuCallbackSettings = "menu:settings"
	MenuCallbackRefresh  = "menu:refresh"
)

func handleMenuCmd(ctx *ext.Context, u *ext.Update) error {
	return showMainMenu(ctx, u.GetUserChat().GetID())
}

func showMainMenu(ctx *ext.Context, chatID int64, msgID ...int) error {
	// Get running tasks count
	runningTasks := core.GetRunningTasks(context.Background())
	queuedTasks := core.GetQueuedTasks(context.Background())

	statusText := "📊 *SaveAny-Bot* - 主菜单\n\n"
	statusText += "━━━━━━━━━━━━━━\n"
	statusText += fmt.Sprintf("📥 下载中: %d\n", len(runningTasks))
	statusText += fmt.Sprintf("⏳ 队列中: %d\n", len(queuedTasks))
	statusText += fmt.Sprintf("💾 存储: %d\n", len(storage.Storages))
	statusText += "━━━━━━━━━━━━━━\n\n"
	statusText += "选择一个操作:"

	// Build inline keyboard - simple and clean
	markup := &tg.ReplyInlineMarkup{
		Rows: []tg.KeyboardButtonRow{
			{
				Buttons: []tg.KeyboardButtonClass{
					&tg.KeyboardButtonCallback{
						Text: "📊 状态",
						Data: []byte(MenuCallbackStatus),
					},
					&tg.KeyboardButtonCallback{
						Text: "📋 任务",
						Data: []byte(MenuCallbackTasks),
					},
				},
			},
			{
				Buttons: []tg.KeyboardButtonClass{
					&tg.KeyboardButtonCallback{
						Text: "💾 存储位置",
						Data: []byte(MenuCallbackStorages),
					},
					&tg.KeyboardButtonCallback{
						Text: "⚙️ 默认存储",
						Data: []byte(MenuCallbackSettings),
					},
				},
			},
			{
				Buttons: []tg.KeyboardButtonClass{
					&tg.KeyboardButtonCallback{
						Text: "🔄 刷新",
						Data: []byte(MenuCallbackRefresh),
					},
				},
			},
		},
	}

	// Send menu message
	if len(msgID) > 0 {
		_, err := ctx.EditMessage(chatID, &tg.MessagesEditMessageRequest{
			Message:    statusText,
			ID:         msgID[0],
			ReplyMarkup: markup,
		})
		return err
	}

	_, err := ctx.SendMessage(chatID, &tg.MessagesSendMessageRequest{
		Message:    statusText,
		ReplyMarkup: markup,
	})
	return err
}

func handleMenuCallback(ctx *ext.Context, u *ext.Update) error {
	callbackData := string(u.CallbackQuery.Data)
	chatID := u.CallbackQuery.GetUserID()
	msgID := u.CallbackQuery.GetMsgID()

	// Answer callback first
	ctx.AnswerCallback(&tg.MessagesSetBotCallbackAnswerRequest{
		QueryID: u.CallbackQuery.QueryID,
	})

	switch {
	case strings.HasPrefix(callbackData, MenuCallbackStatus):
		return showStatusCallback(ctx, chatID, msgID)
	case strings.HasPrefix(callbackData, MenuCallbackTasks):
		return showTasksCallback(ctx, chatID, msgID)
	case strings.HasPrefix(callbackData, MenuCallbackStorages):
		return showStoragesCallback(ctx, chatID, msgID)
	case strings.HasPrefix(callbackData, MenuCallbackSettings):
		return showSettingsCallback(ctx, chatID, msgID)
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

	statusText := "📊 *系统状态*\n\n"
	statusText += "━━━━━━━━━━━━━━\n"
	statusText += fmt.Sprintf("✅ 状态: 运行中\n")
	statusText += fmt.Sprintf("📥 下载任务: %d\n", len(runningTasks))
	statusText += fmt.Sprintf("⏳ 队列任务: %d\n", len(queuedTasks))
	statusText += fmt.Sprintf("💾 存储数量: %d\n", len(storage.Storages))
	statusText += fmt.Sprintf("⚙️ 工作线程: %d\n", config.C().Workers)
	statusText += fmt.Sprintf("🔄 版本: %s\n", config.Version)
	statusText += "━━━━━━━━━━━━━━"

	// Add back button
	markup := &tg.ReplyInlineMarkup{
		Rows: []tg.KeyboardButtonRow{
			{
				Buttons: []tg.KeyboardButtonClass{
					&tg.KeyboardButtonCallback{
						Text: "🔙 返回菜单",
						Data: []byte(MenuCallbackRefresh),
					},
				},
			},
		},
	}

	_, err := ctx.EditMessage(chatID, &tg.MessagesEditMessageRequest{
		Message:    statusText,
		ID:         msgID,
		ReplyMarkup: markup,
	})
	return err
}

func showTasksCallback(ctx *ext.Context, chatID int64, msgID int) error {
	runningTasks := core.GetRunningTasks(context.Background())
	queuedTasks := core.GetQueuedTasks(context.Background())

	var tasksText string

	if len(runningTasks) == 0 && len(queuedTasks) == 0 {
		tasksText = "📋 *任务列表*\n\n暂无任务"
	} else {
		tasksText = "📋 *任务列表*\n\n"

		if len(runningTasks) > 0 {
			tasksText += "📥 *下载中:*\n"
			for i, task := range runningTasks {
				if i >= 5 {
					tasksText += fmt.Sprintf("\n... 还有 %d 个", len(runningTasks)-5)
					break
				}
				tasksText += fmt.Sprintf("• %s\n", task.Title)
			}
		}

		if len(queuedTasks) > 0 {
			tasksText += "\n⏳ *队列中:*\n"
			for i, task := range queuedTasks {
				if i >= 5 {
					tasksText += fmt.Sprintf("\n... 还有 %d 个", len(queuedTasks)-5)
					break
				}
				tasksText += fmt.Sprintf("• %s\n", task.Title)
			}
		}
	}

	// Add action buttons
	markup := &tg.ReplyInlineMarkup{
		Rows: []tg.KeyboardButtonRow{
			{
				Buttons: []tg.KeyboardButtonClass{
					&tg.KeyboardButtonCallback{
						Text: "🔄 刷新",
						Data: []byte(MenuCallbackTasks),
					},
					&tg.KeyboardButtonCallback{
						Text: "🔙 返回",
						Data: []byte(MenuCallbackRefresh),
					},
				},
			},
		},
	}

	_, err := ctx.EditMessage(chatID, &tg.MessagesEditMessageRequest{
		Message:    tasksText,
		ID:         msgID,
		ReplyMarkup: markup,
	})
	return err
}

func showStoragesCallback(ctx *ext.Context, chatID int64, msgID int) error {
	storagesText := "💾 *存储位置*\n\n"

	for name, s := range storage.Storages {
		storType := s.Type().String()
		storagesText += fmt.Sprintf("• *%s* (%s)\n", name, storType)
	}

	if len(storage.Storages) == 0 {
		storagesText += "_暂无存储配置_"
	}

	// Add back button
	markup := &tg.ReplyInlineMarkup{
		Rows: []tg.KeyboardButtonRow{
			{
				Buttons: []tg.KeyboardButtonClass{
					&tg.KeyboardButtonCallback{
						Text: "🔙 返回菜单",
						Data: []byte(MenuCallbackRefresh),
					},
				},
			},
		},
	}

	_, err := ctx.EditMessage(chatID, &tg.MessagesEditMessageRequest{
		Message:    storagesText,
		ID:         msgID,
		ReplyMarkup: markup,
	})
	return err
}

func showSettingsCallback(ctx *ext.Context, chatID int64, msgID int) error {
	settingsText := "⚙️ *设置*\n\n"
	settingsText += "使用 /silent 命令设置默认存储位置\n\n"
	settingsText += "设置后，发送文件将自动保存到默认位置，无需每次选择。\n\n"
	settingsText += "_发送 /silent 来选择默认存储_"

	// Add back button
	markup := &tg.ReplyInlineMarkup{
		Rows: []tg.KeyboardButtonRow{
			{
				Buttons: []tg.KeyboardButtonClass{
					&tg.KeyboardButtonCallback{
						Text: "🔙 返回菜单",
						Data: []byte(MenuCallbackRefresh),
					},
				},
			},
		},
	}

	_, err := ctx.EditMessage(chatID, &tg.MessagesEditMessageRequest{
		Message:    settingsText,
		ID:         msgID,
		ReplyMarkup: markup,
	})
	return err
}
