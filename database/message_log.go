package database

import (
	"context"
	"encoding/json"
	"time"
)

// LogMessage stores a Telegram message for debugging
func LogMessage(ctx context.Context, chatID, userID int64, messageType, message, rawData string) error {
	msg := MessageLog{
		ChatID:     chatID,
		UserID:     userID,
		Message:    message,
		MessageType: messageType,
		RawData:    rawData,
	}
	return db.WithContext(ctx).Create(&msg).Error
}

// GetMessageLogs retrieves recent message logs
func GetMessageLogs(ctx context.Context, limit int) ([]MessageLog, error) {
	var logs []MessageLog
	err := db.WithContext(ctx).
		Order("created_at DESC").
		Limit(limit).
		Find(&logs).Error
	return logs, err
}

// GetMessageLogsByUser retrieves messages for a specific user
func GetMessageLogsByUser(ctx context.Context, userID int64, limit int) ([]MessageLog, error) {
	var logs []MessageLog
	err := db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Find(&logs).Error
	return logs, err
}

// GetMessageLogsByChat retrieves messages from a specific chat
func GetMessageLogsByChat(ctx context.Context, chatID int64, limit int) ([]MessageLog, error) {
	var logs []MessageLog
	err := db.WithContext(ctx).
		Where("chat_id = ?", chatID).
		Order("created_at DESC").
		Limit(limit).
		Find(&logs).Error
	return logs, err
}

// ClearMessageLogs clears all message logs
func ClearMessageLogs(ctx context.Context) error {
	return db.WithContext(ctx).Unscoped().Delete(&MessageLog{}).Error
}

// GetMessageStats returns message statistics
func GetMessageStats(ctx context.Context) (map[string]int64, error) {
	var total int64
	var today int64
	
	// Total count
	db.WithContext(ctx).Model(&MessageLog{}).Count(&total)
	
	// Today's count
	startOfDay := time.Now().Truncate(24 * time.Hour)
	db.WithContext(ctx).Model(&MessageLog{}).Where("created_at >= ?", startOfDay).Count(&today)
	
	// Count by type
	typeCounts := make(map[string]int64)
	var results []struct {
		MessageType string
		Count       int64
	}
	db.WithContext(ctx).Model(&MessageLog{}).
		Select("message_type, COUNT(*) as count").
		Group("message_type").
		Scan(&results)
	
	for _, r := range results {
		typeCounts[r.MessageType] = r.Count
	}
	
	stats := map[string]int64{
		"total": total,
		"today": today,
	}
	for k, v := range typeCounts {
		stats[k] = v
	}
	
	return stats, nil
}

// MessageToJSON converts a message update to JSON string for storage
func MessageToJSON(v interface{}) string {
	data, _ := json.Marshal(v)
	return string(data)
}
