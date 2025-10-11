package constants

import "sync"

const (
	IS_PRODUCTION        bool = false
	PLAYER_POOL_CAP      int  = 100
	DIGIT_LOWER_LIMIT    int  = 1
	DIGIT_UPPER_LIMIT    int  = 8
	MAX_ROWS_DISPLAYED   int  = 10
	LEADERBOARD_MAX_ROWS int  = 100 // not used for now
)

type AutoInc struct {
	lock *sync.Mutex
	id   int
}

func NewAutoInc() *AutoInc {
	return &AutoInc{
		&sync.Mutex{},
		0,
	}
}

// Function for generating auto incrementing id.
func (a *AutoInc) ID() (id int) {
	a.lock.Lock()
	defer a.lock.Unlock()

	id = a.id
	a.id++
	return
}
