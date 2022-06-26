package goflakes

import (
	"fmt"
	"sync"
	"time"
)

type SnowflakeGenerator struct {
	generatedCount int64
	generatedMutex sync.Mutex
	instance       int64
	epoch          time.Time
}

type AsyncReturn struct {
	id  int64
	err error
}

const (
	idLength                = 63
	timestampLength         = 41
	instanceLength          = 10
	sequenceLength          = 12
	timestapShift           = idLength - timestampLength
	instanceShift           = timestapShift - instanceLength
	latestStorableTime      = (1 << timestampLength) - 1
	biggestStorableInstance = (1 << instanceLength) - 1
	biggestStorableSequence = (1 << sequenceLength) - 1
	resetSequence           = 1 << sequenceLength
)

func NewSnowflakeGenerator(epoch time.Time, instance int64) (SnowflakeGenerator, error) {
	if epoch.UnixMilli() > latestStorableTime {
		return SnowflakeGenerator{}, fmt.Errorf("%v (epoch) is to late to fit into any ids.", epoch)
	}
	if instance > biggestStorableInstance {
		return SnowflakeGenerator{}, fmt.Errorf("%v (nodeid) is to big to fit into any ids.", instance)
	}

	return SnowflakeGenerator{
		generatedCount: 0,
		generatedMutex: sync.Mutex{},
		instance:       instance,
		epoch:          epoch,
	}, nil
}

func (s *SnowflakeGenerator) Generate() (int64, error) {
	now := time.Now()
	timestamp := now.Sub(s.epoch).Milliseconds()
	if timestamp < 0 {
		return 0, fmt.Errorf("Now (%v) is before epoch (%v)", now, s.epoch)
	}
	if timestamp > (1<<41 - 1) {
		return 0, fmt.Errorf("41-bit Integer overflow for timestamp. ")
	}
	s.generatedMutex.Lock()
	s.generatedCount++
	var id int64 = timestamp<<timestapShift | s.instance<<instanceShift | s.generatedCount%resetSequence
	s.generatedMutex.Unlock()
	return id, nil
}

func (s *SnowflakeGenerator) GenerateMultiple(amount int) ([]int64, error) {
	if amount > biggestStorableSequence {
		return make([]int64, 0), fmt.Errorf("Amount %v is to high, only up to %v ids can be generated at once.", amount, biggestStorableSequence)
	}
	ids := make([]int64, amount)
	for i := 0; i < amount; i++ {
		id, generateError := s.Generate()
		if generateError != nil {
			return ids, generateError
		}
		ids[i] = id
	}
	return ids, nil
}

func (s *SnowflakeGenerator) internalAsyncGenerate(amount int, returnchannel chan AsyncReturn, wg sync.WaitGroup) {
	defer wg.Done()
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
	if amount > biggestStorableSequence {
		return make(chan AsyncReturn, 0), sync.WaitGroup{}, fmt.Errorf("Amount %v is to high, only up to %v ids can be generated at once.", amount, biggestStorableSequence)
	}
	returnchannel := make(chan AsyncReturn, amount)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go s.internalAsyncGenerate(amount, returnchannel, wg)
	return returnchannel, sync.WaitGroup{}, nil
}