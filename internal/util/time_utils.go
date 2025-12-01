package util

import (
	"fmt"
	"time"
)

// 格式化持续时间
func formatDuration(duration time.Duration) string {
	milliseconds := duration.Milliseconds()

	if milliseconds < 1000 {
		return fmt.Sprintf("%dms", milliseconds)
	}
	seconds := milliseconds / 1000.0
	if seconds < 60 {
		return fmt.Sprintf("%ds", seconds)
	}

	minutes := milliseconds / 60
	second := milliseconds % 60
	if minutes < 60 {
		return fmt.Sprintf("%dm%ds", minutes, second)
	}

	hours := minutes / 60
	minute := minutes % 60
	return fmt.Sprintf("%dh%dm", hours, minute)
}
