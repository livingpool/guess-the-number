package service

import (
	"errors"
	"sync"

	"github.com/livingpool/constants"
)

type PlayerPoolRepository interface {
	NewPlayer(answer, timeZone string) *Player
	AddPlayer(player *Player) error
	RemovePlayer(id int) error
	GetPlayer(id int) (*Player, bool)
}

type PlayerPool struct {
	players   map[int]*Player
	lock      *sync.RWMutex
	capacity  int
	autoIncId *constants.AutoInc
}

func NewPlayerPool(capacity int) *PlayerPool {
	return &PlayerPool{
		players:   make(map[int]*Player),
		lock:      new(sync.RWMutex),
		capacity:  capacity,
		autoIncId: constants.NewAutoInc(),
	}
}

func (p *PlayerPool) NewPlayer(answer, timeZone string) *Player {
	return &Player{
		Id:       p.autoIncId.ID(),
		TimeZone: timeZone,
		Answer:   answer,
	}
}

func (p *PlayerPool) AddPlayer(player *Player) error {
	p.lock.Lock()
	defer p.lock.Unlock()

	// If pool is at capacity, return error
	if len(p.players) == p.capacity {
		return errors.New("Player pool at capacity")
	}

	p.players[player.Id] = player
	return nil
}

func (p *PlayerPool) RemovePlayer(id int) error {
	p.lock.Lock()
	defer p.lock.Unlock()

	delete(p.players, id)
	return nil
}

func (p *PlayerPool) GetPlayer(id int) (*Player, bool) {
	p.lock.RLock()
	defer p.lock.RUnlock()

	player, exists := p.players[id]
	return player, exists
}
