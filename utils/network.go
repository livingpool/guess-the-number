package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"strings"
	"time"
)

type IpResponse struct {
	Status   string `json:"status"`
	Message  string `json:"message"`
	TimeZone string `json:"timezone"`
	Query    string `json:"query"`
}

// Note that any security-related use of these headers
// must only be IP addresses added by a trusted proxy.
//
// See https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/X-Forwarded-For.
func ReadUserIP(r *http.Request) net.IP {
	ip := r.Header.Get("X-Real-Ip")
	if res := net.ParseIP(ip); res != nil {
		return res
	}

	for ip := range strings.SplitSeq(r.Header.Get("X-Forwarded-For"), ", ") {
		if res := net.ParseIP(ip); res != nil {
			return res
		}
	}

	if ip, _, err := net.SplitHostPort(r.RemoteAddr); err != nil {
		slog.Error(err.Error())
	} else {
		if res := net.ParseIP(ip); res != nil {
			return res
		}
	}

	return nil
}

// Get geo location info from ip-api.com.
func GetTimeZone(ip net.IP) string {
	defaultZone := "Asia/Taipei"
	if len(ip) == 0 {
		return defaultZone
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	url := fmt.Sprintf("http://ip-api.com/json/%s?fields=status,message,timezone,query", ip.String())
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		slog.Error("setup request error", "err", err.Error())
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		slog.Error("get ip-api.com error", "err", err.Error())
	}
	defer resp.Body.Close()

	result := IpResponse{}
	json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		slog.Error("json decode error", "err", err.Error())
	}

	if result.Status == "fail" {
		slog.Error("get ip-api.com return failed resp", "resp", result)
	}

	if result.TimeZone != "" {
		return result.TimeZone
	}
	return defaultZone
}
