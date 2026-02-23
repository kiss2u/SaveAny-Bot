package handlers

import (
	"errors"

	"github.com/celestix/gotgproto/dispatcher"
	"github.com/celestix/gotgproto/ext"
	"github.com/charmbracelet/log"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/krau/SaveAny-Bot/client/bot/handlers/utils/dirutil"
	"github.com/krau/SaveAny-Bot/common/i18n"
	"github.com/krau/SaveAny-Bot/common/i18n/i18nk"
	"github.com/krau/SaveAny-Bot/config"
	"github.com/krau/SaveAny-Bot/database"
	"github.com/krau/SaveAny-Bot/storage"
)

// Custom error types for better error handling and classification.
// These errors provide semantic meaning for common failure scenarios.
var (
	ErrUnauthorized       = errors.New("unauthorized user")
	ErrUserNotFound       = errors.New("user not found")
	ErrStorageNotFound    = errors.New("storage not found")
	ErrStorageAccess      = errors.New("unable to access storage")
	ErrDefaultStorageNotSet = errors.New("default storage not set")
	ErrDirNotFound        = errors.New("directory not found")
)

// UserError wraps a user-facing error with optional i18n key
type UserError struct {
	Err error
	Key i18nk.Key
}

func (e *UserError) Error() string {
	return e.Err.Error()
}

func (e *UserError) Unwrap() error {
	return e.Err
}

// NewUserError creates a new UserError with an i18n key
func NewUserError(err error, key i18nk.Key) *UserError {
	return &UserError{Err: err, Key: key}
}

// LogError logs an error with context
func LogError(ctx *ext.Context, operation string, err error) {
	logger := log.FromContext(ctx)
	logger.Errorf("Error in %s: %s", operation, err.Error())
}

// HandleHandlerError processes errors from handlers with appropriate user feedback
func HandleHandlerError(ctx *ext.Context, update *ext.Update, err error) error {
	if err == nil {
		return nil
	}

	// Check for user-facing errors first
	var userErr *UserError
	if errors.As(err, &userErr) {
		if userErr.Key != "" {
			ctx.Reply(update, ext.ReplyTextString(i18n.T(userErr.Key, nil)), nil)
		} else {
			ctx.Reply(update, ext.ReplyTextString(i18n.T(i18nk.BotMsgCommonErrorUserGeneric, nil)), nil)
		}
		return dispatcher.EndGroups
	}

	// Log the error for debugging
	LogError(ctx, "handler", err)

	// Send generic error message to user
	ctx.Reply(update, ext.ReplyTextString(i18n.T(i18nk.BotMsgCommonErrorUserGeneric, nil)), nil)
	return dispatcher.EndGroups
}

func checkPermission(ctx *ext.Context, update *ext.Update) error {
	userID := update.GetUserChat().GetID()
	if !slice.Contain(config.C().GetUsersID(), userID) {
		ctx.Reply(update, ext.ReplyTextString(i18n.T(i18nk.BotMsgCommonErrorNoPermission, nil)), nil)
		return dispatcher.EndGroups
	}

	return dispatcher.ContinueGroups
}

func handleSilentMode(next func(*ext.Context, *ext.Update) error, handler func(*ext.Context, *ext.Update) error) func(*ext.Context, *ext.Update) error {
	return func(ctx *ext.Context, update *ext.Update) error {
		userID := update.GetUserChat().GetID()
		user, err := database.GetUserByChatID(ctx, userID)
		if err != nil {
			ctx.Reply(update, ext.ReplyTextString(i18n.T(i18nk.BotMsgCommonErrorGetUserInfoFailed, map[string]any{
				"Error": err.Error(),
			})), nil)
			return dispatcher.EndGroups
		}
		if !user.Silent {
			return next(ctx, update)
		}
		if user.DefaultStorage == "" {
			ctx.Reply(update, ext.ReplyTextString(i18n.T(i18nk.BotMsgCommonErrorDefaultStorageNotSet, nil)), nil)
			return next(ctx, update)
		}
		stor, err := storage.GetStorageByUserIDAndName(ctx, userID, user.DefaultStorage)
		if err != nil {
			ctx.Reply(update, ext.ReplyTextString(i18n.T(i18nk.BotMsgCommonErrorGetStorageFailed, map[string]any{
				"Error": err.Error(),
			})), nil)
			return dispatcher.EndGroups
		}
		if user.DefaultDir != 0 {
			dir, err := database.GetDirByID(ctx, user.DefaultDir)
			if err != nil {
				ctx.Reply(update, ext.ReplyTextString(i18n.T(i18nk.BotMsgCommonErrorGetDirFailed, map[string]any{
					"Error": err.Error(),
				})), nil)
				return next(ctx, update)
			}
			ctx.Context = dirutil.WithContext(ctx.Context, dir)
		}
		ctx.Context = storage.WithContext(ctx.Context, stor)
		return handler(ctx, update)
	}
}
