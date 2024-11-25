package pkg

import (
	"fmt"
)

const (
	_   = iota             // Ignore first value by assigning to blank identifier
	KiB = 1 << (10 * iota) // 1024 bytes
	MiB                    // 1024 * 1024 bytes
	GiB                    // 1024 * 1024 * 1024 bytes
)

func FormatBytes(bytes int64) string {
	switch {
	case bytes >= GiB:
		return fmt.Sprintf("%.2fGiB", float64(bytes)/GiB)
	case bytes >= MiB:
		return fmt.Sprintf("%.2fMiB", float64(bytes)/MiB)
	case bytes >= KiB:
		return fmt.Sprintf("%.2fKiB", float64(bytes)/KiB)
	default:
		return fmt.Sprintf("%dB", bytes)
	}
}
