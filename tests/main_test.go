package goflakes_test

import (
	"math"
	"testing"
	"time"

	"github.com/blocklisted/goflakes"
	"github.com/blocklisted/goflakes/mock"
)

var TimeEpoch = time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)
var SnowflakeGenerator, snowerror = goflakes.NewSnowflakeGenerator(TimeEpoch, 1<<9)

func TestComputeAndExtractFunction(t *testing.T) {
	var timestamp int64 = 3600000
	var instance int64 = 394
	var sequence int64 = 874
	var correct_id int64 = 15099496014698

	var computed_id = goflakes.ComputeID(timestamp, instance, sequence)
	var extracted_timestamp, extracted_instance, extracted_sequence = goflakes.ExtractId(computed_id)
	if computed_id != correct_id {
		t.Errorf(`computed_id %v was not equal to correct_id %v
input values were               %v, %v, %v
output from extraction function %v, %v, %v`, computed_id, correct_id, timestamp, instance, sequence, extracted_timestamp, extracted_instance, extracted_sequence)
	}
	if extracted_timestamp != timestamp || extracted_instance != instance || extracted_sequence != sequence {
		t.Errorf(`extraction function extracted wrong values %v, %v, %v
instead of %v, %v, %v`, extracted_timestamp, extracted_instance, extracted_sequence, timestamp, instance, sequence)
	}
}

func TestGenerate(t *testing.T) {
	currenttime := time.Date(2022, 6, 9, 6, 9, 6, 9, time.UTC)
	mock.MockTimeNow = mock.CurryMockTimeNow(currenttime)
	timestamp := currenttime.Sub(TimeEpoch)
	correct_id := goflakes.ComputeID(timestamp.Milliseconds(), 1<<9, 0)
	SnowflakeGenerator.ResetGenerated()
	v, f := SnowflakeGenerator.Generate()
	if f != nil {
		t.Errorf("%v, %v, %v", f, TimeEpoch, snowerror)
	}
	if v != correct_id {
		t.Errorf("%v generated not equal to correct id %v", v, correct_id)
	}
	t.Logf("%v", v)
	mock.MockTimeNow = nil
	return
}

func TestNewSnowflake(t *testing.T) {
	_, f := goflakes.NewSnowflakeGenerator(TimeEpoch, 1<<10)
	if f == nil {
		t.Errorf("InstanceID, which is too large can be used.")
	}
}

func TestOverflow(t *testing.T) {
	SnowflakeGenerator.ResetGenerated()
	currenttime := TimeEpoch.Add(time.Duration((1 << 41) * math.Pow10(6)))
	mock.MockTimeNow = mock.CurryMockTimeNow(currenttime)
	_, f := SnowflakeGenerator.Generate()
	if f == nil {
		t.Error("Generate function let 41-bit time over through.")
	}
	mock.MockTimeNow = nil
}

func TestSequence(t *testing.T) {
	SnowflakeGenerator.ResetGenerated()
	for i := 0; i < ((1 << 10) * 2); i++ {
		if f, _ := SnowflakeGenerator.GetNewSequence(TimeEpoch.UnixMilli()); f > ((1 << 12) - 1) {
			t.Errorf("Sequence %v, which is too big (>4095), was generated.", f)
		}

	}
}

func TestAntiDuplicateSystem(t *testing.T) {
	starttime := time.Now()
	// Number is derived from 500 milliseconds * 4096, because there are 4096 sequence available per millisecond.
	mock.MockTimeNow = nil
	values, _ := SnowflakeGenerator.GenerateMultiple(2048000)
	time_taken := time.Now().Sub(starttime)

	visited_map := make(map[int64]bool, 2048000)
	for i, v := range values {
		if visited_map[v] == true {
			ts, _, sq := goflakes.ExtractId(v)
			t.Errorf("Found duplicates in array %v, %v, %v, %v", v, ts, sq, i)
			return
		} else {
			visited_map[v] = true
		}
	}
	if time_taken < (500 * time.Millisecond) {
		t.Error("Throttling is not working")
	}
}

func BenchmarkGenerate(b *testing.B) {
	SnowflakeGenerator.Generate()
}

func BenchmarkAsyncGenerate1000ids(b *testing.B) {
	c, _ := SnowflakeGenerator.AsyncGenerate(1000)
	for {
		_, ok := <-c
		if !ok {
			break
		}
	}
}
