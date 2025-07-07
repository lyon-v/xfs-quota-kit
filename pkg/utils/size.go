package utils

import (
	"fmt"
	"strconv"
	"strings"
)

// ParseSize 解析大小字符串为字节数
func ParseSize(sizeStr string) (uint64, error) {
	sizeStr = strings.ToUpper(strings.TrimSpace(sizeStr))

	var multiplier uint64 = 1
	var numStr string

	switch {
	case strings.HasSuffix(sizeStr, "TB"):
		multiplier = 1024 * 1024 * 1024 * 1024
		numStr = sizeStr[:len(sizeStr)-2]
	case strings.HasSuffix(sizeStr, "GB"):
		multiplier = 1024 * 1024 * 1024
		numStr = sizeStr[:len(sizeStr)-2]
	case strings.HasSuffix(sizeStr, "MB"):
		multiplier = 1024 * 1024
		numStr = sizeStr[:len(sizeStr)-2]
	case strings.HasSuffix(sizeStr, "KB"):
		multiplier = 1024
		numStr = sizeStr[:len(sizeStr)-2]
	case strings.HasSuffix(sizeStr, "B"):
		multiplier = 1
		numStr = sizeStr[:len(sizeStr)-1]
	default:
		numStr = sizeStr
	}

	num, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid size format: %s", sizeStr)
	}

	if num < 0 {
		return 0, fmt.Errorf("size cannot be negative: %s", sizeStr)
	}

	return uint64(num * float64(multiplier)), nil
}

// FormatSize 格式化字节数为可读格式
func FormatSize(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}

	div, exp := uint64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
