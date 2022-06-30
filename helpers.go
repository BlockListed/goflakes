package goflakes

import (
	"fmt"
	"time"

	"github.com/blocklisted/goflakes/constants"
	"github.com/blocklisted/goflakes/mock"
)

// DO NOT USE, public for testing ONLY
// This function will block for up to a millisecond if more than 2^12-1 (4095) sequences are requested in a single millisecond.
// Returns sequence + new timestamp
func (s *SnowflakeGenerator) GetNewSequence(current_unix_time int64) (int64, int64) {
	s.generatedMutex.Lock()
	sequence := s.generatedCount % constants.ResetSequence
	s.generatedCount++
	var new_unix int64
	if sequence == 0 {
		time.Sleep(time.Until(s.LastgeneratedReset.Add(time.Millisecond)))
		new_unix = s.getTimeStamp()
		s.LastgeneratedReset = mock.Time_Now()
	} else {
		new_unix = current_unix_time
	}
	s.generatedMutex.Unlock()
	return sequence, new_unix
}

// No checking
func ComputeID(timestamp int64, instancedid int64, sequence int64) int64 {
	if timestamp > constants.LatestStorableTime || sequence > constants.BiggestStorableSequence {
		fmt.Printf("GOFLAKES - reached unreachable code, because timestamp %v or sequence %v was to big!", timestamp, sequence)
		return 0
	}
	return timestamp<<constants.TimestapShift | instancedid<<constants.InstanceShift | sequence
}

func ExtractId(snowflake int64) (timestamp int64, instanceid int64, sequence int64) {
	timestamp = snowflake & constants.TimestampMask >> constants.TimestapShift
	instanceid = snowflake & constants.InstanceMask >> constants.InstanceShift
	sequence = snowflake & constants.SequenceMask
	return
}

func (s *SnowflakeGenerator) ResetGenerated() {
	s.generatedMutex.Lock()
	s.generatedCount = 0
	s.generatedMutex.Unlock()
}

func (s *SnowflakeGenerator) getTimeStamp() int64 {
	return mock.Time_Now().UnixMilli() - s.epoch
}
