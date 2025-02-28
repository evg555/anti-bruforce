package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidSubnet(t *testing.T) {
	type testCase struct {
		input    string
		expected bool
	}

	testCases := []testCase{
		{"172.0.0.0/24", true},
		{"192.168.1.0/28", true},
		{"10.0.0.0/8", true},
		{"256.0.0.0/24", false}, // Неверный IP (256 недопустимо)
		{"172.0.0.0/33", false}, // Маска больше 32 (для IPv4)
		{"not_a_subnet", false}, // Некорректный ввод
	}

	for _, tc := range testCases {
		t.Run("test subnet: "+tc.input, func(t *testing.T) {
			assert.Equal(t, tc.expected, IsValidSubnet(tc.input))
		})
	}
}

func TestIsValidIpAddress(t *testing.T) {
	type testCase struct {
		input    string
		expected bool
	}

	testCases := []testCase{
		{"172.0.0.1", true},
		{"192.168.1.0", true},
		{"10.0.0.2", true},
		{"0.0.0.0", true},
		{"255.255.255.255", true},
		{"256.0.0.1", false}, // Неверный IP (256 недопустимо)
		{"not_a_ip", false},  // Некорректный ввод
	}

	for _, tc := range testCases {
		t.Run("test ip: "+tc.input, func(t *testing.T) {
			assert.Equal(t, tc.expected, IsValidIpAddress(tc.input))
		})
	}
}
