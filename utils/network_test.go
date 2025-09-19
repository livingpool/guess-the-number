package utils_test

import (
	"net"
	"testing"

	"github.com/livingpool/utils"
	"github.com/stretchr/testify/assert"
)

func TestGetTimeZone(t *testing.T) {
	testcases := []struct {
		ip   string
		want string
	}{
		{"24.48.0.1", "America/Toronto"},
		{"42.70.0.0", "Asia/Taipei"},
	}

	for _, tc := range testcases {
		ip := net.ParseIP(tc.ip)
		timeZone := utils.GetTimeZone(ip)
		assert.Equal(t, tc.want, timeZone)
	}
}
