package goflakes_test

import (
	"testing"
	"time"

	"github.com/blocklisted/goflakes"
	"github.com/blocklisted/goflakes/mock"
)

var TimeEpoch = time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)
var SnowflakeGenerator, snowerror = goflakes.NewSnowflakeGenerator(TimeEpoch, 0b1000000000)

func TestConstantsAndComputeFunction(t *testing.T) {
	return
}

func TestGenerate(t *testing.T) {
	currenttime := time.Date(2022, 6, 9, 6, 9, 6, 9, time.UTC)
	mock.MockTimeNow = mock.CurryMockTimeNow(currenttime)
	SnowflakeGenerator.ResetGenerated()
	v, f := SnowflakeGenerator.Generate()
	if f != nil {
		t.Errorf("%v, %v, %v", f, TimeEpoch, snowerror)
	}
	t.Logf("%v", v)
	return
}

func BenchmarkGenerate(b *testing.B) {
	SnowflakeGenerator.Generate()
}

func BenchmarkAsyncGenerate(b *testing.B) {
	c, w, _ := SnowflakeGenerator.AsyncGenerate(1000)
	w.Wait()
	for {
		_, ok := <-c
		if !ok {
			break
		}
	}
}
