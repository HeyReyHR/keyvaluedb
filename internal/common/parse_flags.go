package common

import (
	"fmt"
	"strconv"
	"strings"
)

func ParseSize(s string) (int, error) {
	s = strings.TrimSpace(strings.ToUpper(s))

	switch {
	case strings.HasSuffix(s, "GB"):
		return parseWithMultiplier(s, "GB", 1024*1024*1024)
	case strings.HasSuffix(s, "MB"):
		return parseWithMultiplier(s, "MB", 1024*1024)
	case strings.HasSuffix(s, "KB"):
		return parseWithMultiplier(s, "KB", 1024)
	case strings.HasSuffix(s, "B"):
		return parseWithMultiplier(s, "B", 1)
	default:
		n, err := strconv.Atoi(s)
		if err != nil {
			return 0, fmt.Errorf("invalid size format %q", s)
		}
		if n < 0 {
			return 0, fmt.Errorf("negative size %q is not allowed", s)
		}
		return n, nil
	}
}

func parseWithMultiplier(s, suffix string, mult int64) (int, error) {
	numPart := strings.TrimSpace(strings.TrimSuffix(s, suffix))
	n, err := strconv.ParseInt(numPart, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid size value %q: %w", s, err)
	}
	if n < 0 {
		return 0, fmt.Errorf("negative size %q is not allowed", s)
	}
	return int(n * mult), nil
}
