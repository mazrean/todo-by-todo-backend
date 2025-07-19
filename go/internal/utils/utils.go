package utils

import (
	"fmt"
	"strconv"
)

func ParseInt64(s string) (int64, error) {
	num, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("stringをint64に変換できませんでした: %w", err)
	}
	return num, nil
}
