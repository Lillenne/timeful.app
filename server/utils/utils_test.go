package utils

import (
	"os"
	"testing"
)

func TestIsSelfHostedPremiumEnabled(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		expected bool
	}{
		{"Empty value", "", false},
		{"True value", "true", true},
		{"1 value", "1", true},
		{"Yes value", "yes", true},
		{"False value", "false", false},
		{"0 value", "0", false},
		{"No value", "no", false},
		{"Random value", "random", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variable
			if tt.envValue != "" {
				os.Setenv("SELF_HOSTED_PREMIUM", tt.envValue)
			} else {
				os.Unsetenv("SELF_HOSTED_PREMIUM")
			}

			// Test the function
			result := IsSelfHostedPremiumEnabled()
			if result != tt.expected {
				t.Errorf("IsSelfHostedPremiumEnabled() = %v, want %v (with SELF_HOSTED_PREMIUM=%q)", result, tt.expected, tt.envValue)
			}
		})
	}
}
