package utils

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"math"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/livingpool/constants"
)

type ValiateCaptchaResp struct {
	Success    bool     `json:"success"`
	ErrorCodes []string `json:"error-codes"`
}

const maxRetries = 3

func ValidateCaptcha(ctx context.Context, token string) error {
	url := "https://challenges.cloudflare.com/turnstile/v0/siteverify"

	secret := os.Getenv("CLOUDFLARE_TURNSTILE_SECRET_KEY")
	if !constants.IS_PRODUCTION || secret == "" {
		secret = "1x0000000000000000000000000000000AA"
	}

	idempotencyKey := uuid.NewString()
	res := &ValiateCaptchaResp{}

	for i := range maxRetries {
		reqData := strings.NewReader(
			fmt.Sprintf(`{"secret":"%s","response":"%s","idempotency_key":"%s"}`,
				secret,
				token,
				idempotencyKey,
			))
		req, err := http.NewRequestWithContext(ctx, "POST", url, reqData)
		if err != nil {
			return fmt.Errorf("failed to construct request: %w", err)
		}

		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{Timeout: 5 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("failed to send request: %w", err)
		}
		defer resp.Body.Close()

		respData, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body: %w", err)
		}

		if err := json.Unmarshal(respData, res); err != nil {
			return fmt.Errorf("failed to unmarshal captcha resp: %w", err)
		}

		if res.Success {
			return nil
		}

		// exponential backoff
		time.Sleep(time.Duration(math.Pow(2, float64(i))) * time.Second)
	}

	slog.Info("captcha validation failed", "success", res.Success, "error_codes", res.ErrorCodes)

	return errors.New(res.ErrorCodes[0])
}
