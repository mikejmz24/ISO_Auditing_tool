package utils

import (
	"fmt"
	"time"
)

// bytesToTime parses time bytes into a time.Time
// Handles multiple time formats including RFC3339
func BytesToTime(b []uint8) (time.Time, error) {
	if b == nil {
		return time.Time{}, fmt.Errorf("nil value cannot be converted to time.Time")
	}

	str := string(b)

	// Try multiple time formats, starting with the one in your error
	formats := []string{
		time.RFC3339,          // "2006-01-02T15:04:05Z07:00" (format in your error)
		"2006-01-02 15:04:05", // Your current format
		"2006-01-02T15:04:05", // ISO8601 without timezone
		"2006-01-02",          // Just date
	}

	var firstErr error
	for _, layout := range formats {
		t, err := time.Parse(layout, str)
		if err == nil {
			return t, nil
		}
		if firstErr == nil {
			firstErr = err
		}
	}

	// If we get here, none of the formats matched
	return time.Time{}, fmt.Errorf("failed to parse time '%s': %w", str, firstErr)
}

// bytesToTimePtr parses time bytes into a *time.Time
// Returns nil if input is nil
func BytesToTimePtr(b []uint8) (*time.Time, error) {
	if b == nil {
		return nil, nil
	}

	t, err := BytesToTime(b)
	if err != nil {
		return nil, err
	}

	return &t, nil
}
