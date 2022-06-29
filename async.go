package goflakes

import (
	"fmt"
	"sync"

	"github.com/blocklisted/goflakes/constants"
)

func (s *SnowflakeGenerator) internalAsyncGenerate(amount int, returnchannel chan AsyncReturn, wg sync.WaitGroup) {
	defer wg.Done()
	defer close(returnchannel)

	for i := 0; i < amount; i++ {
		id, generateError := s.Generate()
		if generateError != nil {
			returnchannel <- AsyncReturn{
				id:  0,
				err: generateError,
			}
		} else {
			returnchannel <- AsyncReturn{
				id:  id,
				err: nil,
			}
		}
	}
}

func (s *SnowflakeGenerator) AsyncGenerate(amount int) (chan AsyncReturn, sync.WaitGroup, error) {
	if amount > constants.BiggestStorableSequence {
		return make(chan AsyncReturn, 0), sync.WaitGroup{}, fmt.Errorf("Amount %v is to high, only up to %v ids can be generated at once.", amount, constants.BiggestStorableSequence)
	}
	returnchannel := make(chan AsyncReturn, amount)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go s.internalAsyncGenerate(amount, returnchannel, wg)
	return returnchannel, sync.WaitGroup{}, nil
}
