// Package handlers provides Telegram bot command handlers and middleware.
package handlers

import (
	"github.com/celestix/gotgproto/ext"
	"github.com/charmbracelet/log"
	"github.com/krau/SaveAny-Bot/common/i18n"
	"github.com/krau/SaveAny-Bot/common/i18n/i18nk"
	"github.com/gotd/td/tg"
)

// ReplyWithError sends a user-friendly error message to the user.
// The error is logged for debugging purposes but only a generic message is shown to the user.
func ReplyWithError(ctx *ext.Context, update *ext.Update, err error, key i18nk.Key) {
	logger := log.FromContext(ctx)
	if err != nil {
		logger.Errorf("Operation failed: %s", err)
	}
	ctx.Reply(update, ext.ReplyTextString(i18n.T(key, nil)), nil)
}

// ReplyWithErrorf sends a user-friendly error message with formatted parameters.
// The error is logged for debugging but only a generic message is shown to the user.
func ReplyWithErrorf(ctx *ext.Context, update *ext.Update, err error, key i18nk.Key, params map[string]any) {
	logger := log.FromContext(ctx)
	if err != nil {
		logger.Errorf("Operation failed: %s", err)
	}
	ctx.Reply(update, ext.ReplyTextString(i18n.T(key, params)), nil)
}

// EditMessageWithError edits an existing message with a user-friendly error.
// Used for updating progress messages or inline keyboards with error states.
func EditMessageWithError(ctx *ext.Context, chatID int64, msgID int, err error, key i18nk.Key) {
	logger := log.FromContext(ctx)
	if err != nil {
		logger.Errorf("Operation failed: %s", err)
	}
	ctx.EditMessage(chatID, &tg.MessagesEditMessageRequest{
		ID:      msgID,
		Message: i18n.T(key, nil),
	})
}

// EditMessageWithErrorf edits an existing message with a user-friendly error and params.
// Used for updating progress messages with formatted error information.
func EditMessageWithErrorf(ctx *ext.Context, chatID int64, msgID int, err error, key i18nk.Key, params map[string]any) {
	logger := log.FromContext(ctx)
	if err != nil {
		logger.Errorf("Operation failed: %s", err)
	}
	ctx.EditMessage(chatID, &tg.MessagesEditMessageRequest{
		ID:      msgID,
		Message: i18n.T(key, params),
	})
}

// ReplyWithSuccess sends a success confirmation message to the user.
// Use this for successful operations that don't require detailed feedback.
func ReplyWithSuccess(ctx *ext.Context, update *ext.Update, key i18nk.Key) {
	ctx.Reply(update, ext.ReplyTextString(i18n.T(key, nil)), nil)
}

// ReplyWithSuccessf sends a success message with formatted parameters.
// Use this when you need to include dynamic data in the success message.
func ReplyWithSuccessf(ctx *ext.Context, update *ext.Update, key i18nk.Key, params map[string]any) {
	ctx.Reply(update, ext.ReplyTextString(i18n.T(key, params)), nil)
}

// ReplyWithInfo sends an informational message to the user.
// Use this for status updates, prompts, or general information.
func ReplyWithInfo(ctx *ext.Context, update *ext.Update, key i18nk.Key) {
	ctx.Reply(update, ext.ReplyTextString(i18n.T(key, nil)), nil)
}

// ReplyWithInfof sends an informational message with formatted parameters.
// Use this when you need to include dynamic data in an informational message.
func ReplyWithInfof(ctx *ext.Context, update *ext.Update, key i18nk.Key, params map[string]any) {
	ctx.Reply(update, ext.ReplyTextString(i18n.T(key, params)), nil)
}
