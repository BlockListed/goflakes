package goflakes

import (
	"fmt"
	"sync"
	"time"

	"github.com/blocklisted/goflakes/constants"
)

type SnowflakeGenerator struct {
	generatedCount     int64
	generatedMutex     sync.Mutex
	LastgeneratedReset time.Time
	instance           int64
	epoch              int64
}

type AsyncReturn struct {
	id  int64
	err error
}

func NewSnowflakeGenerator(epoch time.Time, instance int64) (SnowflakeGenerator, error) {
	if instance > constants.BiggestStorableInstance {
		return SnowflakeGenerator{}, fmt.Errorf("%v (nodeid) is to big to fit into any ids.", instance)
	}

	return SnowflakeGenerator{
		generatedCount:     0,
		generatedMutex:     sync.Mutex{},
		LastgeneratedReset: time.Time{},
		instance:           instance,
		epoch:              epoch.UnixMilli(),
	}, nil
}

func (s *SnowflakeGenerator) Generate() (int64, error) {
	timestamp := s.getTimeStamp()
	if timestamp < 0 {
		return 0, fmt.Errorf("Now (%v) is before epoch (%v)", s.epoch+timestamp, s.epoch)
	}
	if timestamp > (1<<41 - 1) {
		return 0, fmt.Errorf("41-bit Integer overflow for timestamp. (%v)", timestamp)
	}
	sequence, ts := s.GetNewSequence(timestamp)
	var id int64 = ComputeID(ts, s.instance, sequence)
	return id, nil
}

func (s *SnowflakeGenerator) GenerateMultiple(amount int) ([]int64, error) {
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
