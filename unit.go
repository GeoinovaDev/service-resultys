package service

import (
	"sync"
	"time"

	"git.resultys.com.br/motor/models/token"
)

// Unit struct
type Unit struct {
	// public
	ID         int          `json:"id"`
	Token      *token.Token `json:"token"`
	Item       interface{}  `json:"data"`
	Status     string       `json:"status"`
	CreateAt   time.Time    `json:"create_at"`
	UpdateAt   time.Time    `json:"update_at"`
	Processing int          `json:"processing"`

	// private
	totalRunning int
	callback     func(*Unit)
	mutex        *sync.Mutex
	wg           *sync.WaitGroup
}

// New cria uma unidade de processamento
func New(tken *token.Token, item interface{}) *Unit {
	unit := &Unit{Token: tken, Item: item}

	unit.mutex = &sync.Mutex{}
	unit.wg = &sync.WaitGroup{}

	return unit
}

// SetStatus ...
func (u *Unit) SetStatus(status string) {
	u.Status = status
	u.UpdateAt = time.Now()
}

// Done callback lançado ao termino do processamento
func (u *Unit) Done(callback func(*Unit)) *Unit {
	u.callback = callback

	return u
}

// Alloc incrementa
func (u *Unit) Alloc(total int) *Unit {
	u.totalRunning = total
	u.wg.Add(total)

	return u
}

// Wait ...
func (u *Unit) Wait() {
	if u.totalRunning > 0 {
		u.wg.Wait()
	}
}

// Release libera uma unidade
func (u *Unit) Release() {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	u.totalRunning--
	if u.totalRunning == 0 {
		if u.callback != nil {
			u.callback(u)
		}

		u.wg.Done()
	}
}
