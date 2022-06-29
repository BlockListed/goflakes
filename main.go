package goflakes

import (
	"fmt"
	"sync"
	"time"

	"github.com/blocklisted/goflakes/constants"
	"github.com/blocklisted/goflakes/mock"
)

type SnowflakeGenerator struct {
	generatedCount     int64
	generatedMutex     sync.Mutex
	LastgeneratedReset int64
	instance           int64
	epoch              time.Time
}

type AsyncReturn struct {
	id  int64
	err error
}

func NewSnowflakeGenerator(epoch time.Time, instance int64) (SnowflakeGenerator, error) {
	if epoch.UnixMilli() > constants.LatestStorableTime {
		return SnowflakeGenerator{}, fmt.Errorf("%v (epoch) is to late to fit into any ids.", epoch)
	}
	if instance > constants.BiggestStorableInstance {
		return SnowflakeGenerator{}, fmt.Errorf("%v (nodeid) is to big to fit into any ids.", instance)
	}

	return SnowflakeGenerator{
		generatedCount:     0,
		generatedMutex:     sync.Mutex{},
		LastgeneratedReset: 0,
		instance:           instance,
		epoch:              epoch,
	}, nil
}

func (s *SnowflakeGenerator) Generate() (int64, error) {
	now := mock.Time_Now()
	timestamp := now.Sub(s.epoch).Milliseconds()
	if timestamp < 0 {
		return 0, fmt.Errorf("Now (%v) is before epoch (%v)", now, s.epoch)
	}
	if timestamp > (1<<41 - 1) {
		return 0, fmt.Errorf("41-bit Integer overflow for timestamp. (%v)", timestamp)
	}
	var id int64 = timestamp<<constants.TimestapShift | s.instance<<constants.InstanceShift | s.getNewSequence()
	return id, nil
}

func (s *SnowflakeGenerator) GenerateMultiple(amount int) ([]int64, error) {
	if amount > constants.BiggestStorableSequence {
		return make([]int64, 0), fmt.Errorf("Amount %v is to high, only up to %v ids can be generated at once.", amount, constants.BiggestStorableSequence)
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
