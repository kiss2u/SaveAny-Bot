package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"slices"

	"github.com/charmbracelet/log"
	"github.com/kiss2u/SaveAny-Bot/client/bot"
	userclient "github.com/kiss2u/SaveAny-Bot/client/user"
	"github.com/kiss2u/SaveAny-Bot/common/cache"
	"github.com/kiss2u/SaveAny-Bot/common/i18n"
	"github.com/kiss2u/SaveAny-Bot/common/notify"
	"github.com/kiss2u/SaveAny-Bot/common/utils/fsutil"
	"github.com/kiss2u/SaveAny-Bot/config"
	"github.com/kiss2u/SaveAny-Bot/core"
	"github.com/kiss2u/SaveAny-Bot/database"
	"github.com/kiss2u/SaveAny-Bot/parsers"
	"github.com/kiss2u/SaveAny-Bot/storage"
	"github.com/kiss2u/SaveAny-Bot/web"
	"github.com/spf13/cobra"
)

func Run(cmd *cobra.Command, _ []string) {
	ctx, cancel := context.WithCancel(cmd.Context())
	logger := log.NewWithOptions(os.Stdout, log.Options{
		Level:           log.DebugLevel,
		ReportTimestamp: true,
		TimeFormat:      time.TimeOnly,
		ReportCaller:    true,
	})
	ctx = log.WithContext(ctx, logger)

	exitChan, err := initAll(ctx, cmd)
	if err != nil {
		logger.Fatal("Init failed", "error", err)
	}
	go func() {
		<-exitChan
		cancel()
	}()

	core.Run(ctx)

	<-ctx.Done()
	logger.Info("Exiting...")
	defer logger.Info("Exit complete")
	cleanCache()
}

func initAll(ctx context.Context, cmd *cobra.Command) (<-chan struct{}, error) {
	configFile := config.GetConfigFile(cmd)
	if err := config.Init(ctx, configFile); err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	cache.Init()
	logger := log.FromContext(ctx)
	i18n.Init(config.C().Lang)
	logger.Info("Initializing...")
	database.Init(ctx)
	// Initialize task state persistence
	database.InitTaskState(ctx)
	storage.LoadStorages(ctx)
	if config.C().Parser.PluginEnable {
		for _, dir := range config.C().Parser.PluginDirs {
			if err := parsers.LoadPlugins(ctx, dir); err != nil {
				logger.Error("Failed to load parser plugins", "dir", dir, "error", err)
			} else {
				logger.Debug("Loaded parser plugins from directory", "dir", dir)
			}
		}
	}
	if config.C().Telegram.Userbot.Enable {
		_, err := userclient.Login(ctx)
		if err != nil {
			logger.Fatal("User login failed", "error", err)
		}
	}
	// Start web server if enabled
	startWebServer(ctx)

	botChan, botClient := bot.Init(ctx)

	// Start health checker with admin notifications
	if botClient != nil {
		healthChecker := bot.NewHealthChecker(botClient, 30*time.Second, 10)

		// Setup admin notifications
		adminIDs := getAdminUserIDs()
		if len(adminIDs) > 0 && bot.ExtContext() != nil {
			adminNotifier := notify.NewAdminNotifier(bot.ExtContext(), adminIDs)
			go adminNotifier.NotifyStartup()

			healthChecker.OnDisconnected = func() {
				go adminNotifier.NotifyDisconnected()
			}
			healthChecker.OnReconnected = func() {
				go adminNotifier.NotifyReconnected()
			}
			healthChecker.OnReconnectFailed = func() {
				go adminNotifier.NotifyReconnectFailed()
			}
		}

		go healthChecker.Start(ctx)
		log.Info("Health checker started")
	}

	return botChan, nil
}

func getAdminUserIDs() []int64 {
	var ids []int64
	for _, user := range config.C().Users {
		ids = append(ids, user.ID)
	}
	return ids
}

func startWebServer(ctx context.Context) *web.Server {
	webConfig := config.C().Web
	if !webConfig.Enable {
		return nil
	}

	server := web.New(ctx, &webConfig)
	go func() {
		if err := server.Run(); err != nil {
			log.FromContext(ctx).Error("Web server failed", "error", err)
		}
	}()

	log.FromContext(ctx).Info("Web server started", "host", webConfig.Host, "port", webConfig.Port)
	return server
}

func cleanCache() {
	if config.C().NoCleanCache {
		return
	}
	if config.C().Temp.BasePath != "" && !config.C().Stream {
		if slices.Contains([]string{"/", ".", "\\", ".."}, filepath.Clean(config.C().Temp.BasePath)) {
			log.Error("Invalid cache directory", "path", config.C().Temp.BasePath)
			return
		}
		currentDir, err := os.Getwd()
		if err != nil {
			log.Error("Failed to get working directory", "error", err)
			return
		}
		cachePath := filepath.Join(currentDir, config.C().Temp.BasePath)
		cachePath, err = filepath.Abs(cachePath)
		if err != nil {
			log.Error("Failed to get absolute cache path", "error", err)
			return
		}
		log.Info("Cleaning cache directory", "path", cachePath)
		if err := fsutil.RemoveAllInDir(cachePath); err != nil {
			log.Error("Failed to clean cache directory", "error", err)
		}
	}
}
