package service

import (
	"log/slog"
	"time"
)

type TimeProviderRepository interface {
	Now(timeZone string) time.Time
}

type RealTimeProvider struct{}

func (rtp *RealTimeProvider) Now(timeZone string) time.Time {
	loc, err := time.LoadLocation(timeZone)
	if err != nil {
		slog.Error("invalid time zone", "timeZone", timeZone)
		loc, _ = time.LoadLocation("Asia/Taipei")
	}
	return time.Now().In(loc)
}

type Player struct {
	Id           int
	TimeZone     string
	Answer       string
	GuessResults GuessResults
}

type FormData struct {
	Digit    int
	Start    string
	End      string
	Error    string
	PlayerId int
}

type ResultRow struct {
	TimeStamp string
	Guess     string
	Result    string
}

type GuessResults struct {
	Rows []ResultRow
}
