package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/livingpool/constants"
	"github.com/livingpool/middleware"
	"github.com/livingpool/service"
	"github.com/livingpool/utils"
	"github.com/livingpool/views"
)

type LeaderboardHandler struct {
	renderer    views.TemplatesRepository
	leaderboard service.LeaderboardRepository
}

func NewLeaderboardHandler(r views.TemplatesRepository, l service.LeaderboardRepository) *LeaderboardHandler {
	return &LeaderboardHandler{
		renderer:    r,
		leaderboard: l,
	}
}

func (h *LeaderboardHandler) SaveRecord(w http.ResponseWriter, r *http.Request) {
	reqId := r.Context().Value(middleware.RequestIdKey).(string)

	token := r.URL.Query().Get("token")
	if token == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("missing token in query param"))
		slog.Warn("save record request did not carry an captcha token")
		return
	}

	captchaCtx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	if err := utils.ValidateCaptcha(captchaCtx, token); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(fmt.Appendf([]byte{}, "failed to validate captcha: %v", err))
		slog.Error("validate captcha failed", "token", token, "err", err.Error())
		return
	}

	var record service.Record

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(&record); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("failed to decode request body"))
		slog.Error("decode json failed", "reqId", reqId, "err", err.Error())
		return
	}

	if record.Digits < constants.DIGIT_LOWER_LIMIT || record.Digits > constants.DIGIT_UPPER_LIMIT {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("%d is out of range", record.Digits)))
		return
	}

	record.Name = strings.TrimSpace(record.Name)
	if len(record.Name) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("name too short"))
		return
	}

	if record.Attempts < 1 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("attempts cannot < 1"))
		return
	}

	if err := h.leaderboard.Insert(r.Context(), record); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		slog.Error("insert leaderboard", "reqId", reqId, "err", err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("record inserted"))
}

func (h *LeaderboardHandler) ShowLeaderboard(w http.ResponseWriter, r *http.Request) {
	reqId := r.Context().Value(middleware.RequestIdKey).(string)
	digit := r.URL.Query().Get("digit")
	name := strings.TrimSpace(r.URL.Query().Get("name"))

	boardId, err := strconv.Atoi(digit)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("%s is not a integer", digit)))
		return
	}

	if boardId < constants.DIGIT_LOWER_LIMIT || boardId > constants.DIGIT_UPPER_LIMIT {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("%d is out of range", boardId)))
		return
	}

	result, err := h.leaderboard.Get(r.Context(), boardId, name)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		slog.Error("get leaderboard failed", "reqId", reqId, "boardId", boardId, "name", name, "err", err.Error())
		w.Write([]byte("record not found"))
		return
	}

	if err := h.renderer.Render(w, "leaderboard", result); err != nil {
		slog.Error("render leaderboard error", "err", err.Error())
	}
}
