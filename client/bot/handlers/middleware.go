package handlers

import (
	"github.com/celestix/gotgproto/dispatcher"
	"github.com/celestix/gotgproto/ext"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/kiss2u/SaveAny-Bot/client/bot/handlers/utils/dirutil"
	"github.com/kiss2u/SaveAny-Bot/common/i18n"
	"github.com/kiss2u/SaveAny-Bot/common/i18n/i18nk"
	"github.com/kiss2u/SaveAny-Bot/config"
	"github.com/kiss2u/SaveAny-Bot/database"
	"github.com/kiss2u/SaveAny-Bot/storage"
)

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
