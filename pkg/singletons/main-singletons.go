package singletons

import (
	sync "sync"

	validation "github.com/beckend/go-config/pkg/validation"
)

// Singletons for the whole app
type Singletons struct {
	Validation validation.Validator
}

var (
	instance *Singletons
	doOnce   sync.Once
)

// New returns struct Singletons
func New() *Singletons {
	doOnce.Do(func() {
		instance = &Singletons{}
		instance.Validation = validation.New()
	})

	return instance
}
