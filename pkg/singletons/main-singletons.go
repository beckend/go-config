package singletons

import (
	sync "sync"

	validation "github.com/beckend/go-config/pkg/validation"
)

// Singletons for the whole app
type Singletons struct {
	Validation validation.GetValidatorReturn
}

var (
	instance *Singletons
	doOnce   sync.Once
)

// GetSingletons returns struct Singletons
func GetSingletons() *Singletons {
	doOnce.Do(func() {
		instance = &Singletons{}
		instance.Validation = validation.GetValidator()
	})

	return instance
}
