package service

import (
	"context"
	"database/sql"
	"fmt"
	"math"

	_ "github.com/tursodatabase/libsql-client-go/libsql"

	"github.com/livingpool/constants"
)

type LeaderboardInterface interface {
	Insert(ctx context.Context, data Record) error
	Get(ctx context.Context, boardId int, name string) ([]Record, error)
	Close() error
}

type Leaderboard struct {
	DB *sql.DB
}

type Record struct {
	Digits   int    `json:"digits"`
	Name     string `json:"name"`
	Attempts int    `json:"attempts"`
}

func NewLeaderboard(db *sql.DB) *Leaderboard {
	return &Leaderboard{db}
}

// inserts a record in the corresponding board if it doesn't already exists, otherwise updates it
func (lb *Leaderboard) Insert(ctx context.Context, data Record) error {
	var minAttempt int
	if err := lb.DB.QueryRowContext(ctx, fmt.Sprintf("select min(attempts) from board%d where name = '%s'", data.Digits, data.Name)).Scan(&minAttempt); err != nil {
		minAttempt = math.MaxInt64
	}

	tx, err := lb.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	stmt, err := tx.PrepareContext(ctx, fmt.Sprintf("insert or replace into board%d (name, attempts) values(?, ?)", data.Digits))
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err = stmt.ExecContext(ctx, data.Name, min(minAttempt, data.Attempts)); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

// if name is non-null, result must contain name
func (lb *Leaderboard) Get(ctx context.Context, boardId int, name string) ([]Record, error) {
	if boardId < constants.DIGIT_LOWER_LIMIT || boardId > constants.DIGIT_UPPER_LIMIT {
		return nil, fmt.Errorf("board id out of range")
	}

	rows, err := lb.DB.QueryContext(ctx, fmt.Sprintf("select name, attempts from board%d order by attempts, name limit %d", boardId, constants.MAX_ROWS_DISPLAYED))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]Record, 0, constants.MAX_ROWS_DISPLAYED)
	nameIncluded := false

	for rows.Next() {
		record := Record{Digits: boardId}
		if err := rows.Scan(&record.Name, &record.Attempts); err != nil {
			return nil, err
		}
		if record.Name == name {
			nameIncluded = true
		}
		result = append(result, record)
	}

	// replace the last row with the player's name if applicable
	if name != "" && !nameIncluded {
		record := Record{Digits: boardId}
		row := lb.DB.QueryRowContext(ctx, fmt.Sprintf("select name, attempts from board%d where name = '%s'", boardId, name))
		if err := row.Scan(&record.Name, &record.Attempts); err != nil {
			return nil, err
		}
		result[constants.MAX_ROWS_DISPLAYED-1] = record
	}

	return result, nil
}

// there is generally no need to close the db connection
func (lb *Leaderboard) Close() error {
	return lb.DB.Close()
}
