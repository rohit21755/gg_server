package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

// GenerateRandomString generates a random string of specified length
func GenerateRandomString(length int) (string, error) {
	bytes := make([]byte, length/2)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// GenerateReferralCode generates a unique referral code
func GenerateReferralCode() (string, error) {
	code, err := GenerateRandomString(8)
	if err != nil {
		return "", err
	}
	return strings.ToUpper(code), nil
}

// FormatDuration formats a duration to human-readable string
func FormatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%ds", int(d.Seconds()))
	}
	if d < time.Hour {
		return fmt.Sprintf("%dm", int(d.Minutes()))
	}
	if d < 24*time.Hour {
		return fmt.Sprintf("%dh", int(d.Hours()))
	}
	return fmt.Sprintf("%dd", int(d.Hours()/24))
}

// TruncateString truncates a string to specified length
func TruncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// Contains checks if a string slice contains a value
func Contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// PointerToString converts a string pointer to string
func PointerToString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// StringToPointer converts a string to string pointer
func StringToPointer(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// IntToPointer converts an int to int pointer
func IntToPointer(i int) *int {
	return &i
}

// PointerToInt converts an int pointer to int
func PointerToInt(i *int) int {
	if i == nil {
		return 0
	}
	return *i
}

