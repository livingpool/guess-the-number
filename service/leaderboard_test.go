package service_test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/mattn/go-sqlite3"

	"github.com/livingpool/constants"
	"github.com/livingpool/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// Note: package suite does not support parallel tests

type TestSuite struct {
	suite.Suite
	db *sql.DB
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (ts *TestSuite) SetupSuite() {
	db, err := sql.Open("sqlite3", "file:test.db?mode=memory&cache=shared")
	if err != nil {
		ts.Suite.T().Errorf("could not connect to database: %v", err)
	}
	ts.db = db

	for i := constants.DIGIT_LOWER_LIMIT; i <= constants.DIGIT_UPPER_LIMIT; i++ {
		if _, err := ts.db.Exec(fmt.Sprintf("create table board%d (id integer primary key, name text unique, attempts integer)", i)); err != nil {
			ts.Suite.T().Errorf("create table board%d failed: %v", i, err)
		}
	}
}

// before each
func (ts *TestSuite) SetupTest() {
	for i := constants.DIGIT_LOWER_LIMIT; i <= constants.DIGIT_UPPER_LIMIT; i++ {
		for j := range constants.MAX_ROWS_DISPLAYED + 1 {
			if _, err := ts.db.Exec(fmt.Sprintf("insert into board%d (name, attempts) values ('tim%[2]d', %[2]d)", i, j+1)); err != nil {
				ts.Suite.T().Errorf("insert into board%d failed: %v", i, err)
			}
		}
	}
}

// after each
func (ts *TestSuite) TearDownTest() {
	for i := constants.DIGIT_LOWER_LIMIT; i <= constants.DIGIT_UPPER_LIMIT; i++ {
		if _, err := ts.db.Exec(fmt.Sprintf("delete from board%d", i)); err != nil {
			ts.Suite.T().Errorf("delete from board%d failed: %v", i, err)
		}
	}
}

func (ts *TestSuite) TearDownSuite() {
	ts.db.Close()
}

func (ts *TestSuite) TestInsert() {
	record := service.Record{
		Digits:   5,
		Name:     "Fuyen Guo",
		Attempts: 12,
	}

	leaderboard := service.Leaderboard{DB: ts.db}

	err := leaderboard.Insert(context.TODO(), record)
	require.NoError(ts.Suite.T(), err)

	inserted := service.Record{Digits: 5}
	err = ts.db.QueryRow(fmt.Sprintf("select name, attempts from board%d where name = '%s'", record.Digits, record.Name)).Scan(&inserted.Name, &inserted.Attempts)
	require.NoError(ts.Suite.T(), err)
	assert.Equal(ts.Suite.T(), record, inserted)
}

func (ts *TestSuite) TestGet() {
	testcases := []struct {
		record service.Record
		result bool
	}{
		{service.Record{3, "tim3", 3}, true},
		{service.Record{3, "", 3}, true},
		{service.Record{5, "Fuyen Guo", 12}, false},
	}

	leaderboard := service.Leaderboard{DB: ts.db}
	for _, tc := range testcases {
		obj, err := leaderboard.Get(context.TODO(), tc.record.Digits, tc.record.Name)
		if tc.result {
			require.NoError(ts.Suite.T(), err)
			assert.Equal(ts.Suite.T(), len(obj), constants.MAX_ROWS_DISPLAYED)
			if tc.record.Name != "" {
				assert.Contains(ts.Suite.T(), obj, tc.record)
			}
		} else {
			assert.ErrorContains(ts.Suite.T(), err, "no rows")
		}
	}
}
