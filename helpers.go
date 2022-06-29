package goflakes

import (
	"time"

	"github.com/blocklisted/goflakes/constants"
	"github.com/blocklisted/goflakes/mock"
)

// This function will block for up to a millisecond if more than 2^12-1 (4095) sequences are requested in a single millisecond.
func (s *SnowflakeGenerator) getNewSequence() int64 {
	s.generatedMutex.Lock()
	sequence := s.generatedCount % constants.ResetSequence
	s.generatedCount++
	if sequence == 0 {
		now_unix := mock.Time_Now().UnixMilli()
		if s.LastgeneratedReset == now_unix {
			time.Sleep(time.Until(time.UnixMilli(now_unix + 1)))
		}
		s.LastgeneratedReset = now_unix
	}
	s.generatedMutex.Unlock()
	return sequence
}

// No checking
func ComputeID(timestamp int64, instancedid int64, sequence int64) int64 {
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
