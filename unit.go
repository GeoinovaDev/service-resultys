package service

import (
	"sync"

	"git.resultys.com.br/lib/lower/promise"

	"git.resultys.com.br/motor/models/token"
)

// Unit struct
type Unit struct {
	// public
	Token  *token.Token
	Item   interface{}
	Finish *promise.Promise

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
	unit.Finish = promise.New()

	return unit
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
