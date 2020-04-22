package service

import "time"

// Service interface
type Service interface {
	Add(*Unit)
	Load()
	Reload()
	Stats() time.Duration
}
