package pdf

import (
	"fmt"
	"strconv"
	"strings"
)

func parseMargin(value string) (float64, error) {
	value = strings.TrimSpace(value)

	switch {
	case strings.HasSuffix(value, "mm"):
		v, err := strconv.ParseFloat(strings.TrimSuffix(value, "mm"), 64)
		if err != nil {
			return 0, err
		}
		return v / 25.4, nil // mm → inch

	case strings.HasSuffix(value, "cm"):
		v, err := strconv.ParseFloat(strings.TrimSuffix(value, "cm"), 64)
		if err != nil {
			return 0, err
		}
		return v / 2.54, nil // cm → inch

	case strings.HasSuffix(value, "in"):
		return strconv.ParseFloat(strings.TrimSuffix(value, "in"), 64)

	default:
		return 0, fmt.Errorf("invalid margin unit: %s", value)
	}
}
