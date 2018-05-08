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
}

// New cria uma unidade de processamento
func New(token *token.Token) *Unit {
	unit := &Unit{Token: token}

	unit.mutex = &sync.Mutex{}
	unit.Finish = &promise.Promise{}

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

	return u
}

// Release libera uma unidade
func (u *Unit) Release() {
	u.mutex.Lock()

	u.totalRunning--
	if u.totalRunning == 0 && u.callback != nil {
		u.callback(u)
	}

	u.mutex.Unlock()
}
