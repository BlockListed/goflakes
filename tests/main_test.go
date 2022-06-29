package goflakes_test

import (
	"testing"
	"time"

	"github.com/blocklisted/goflakes"
)

var TimeEpoch, parseerror = time.Parse(time.RFC3339, "2022-01-01T0:00:00Z")
var SnowflakeGenerator, snowerror = goflakes.NewSnowflakeGenerator(TimeEpoch, 0b1000000000)

func TestGenerate(t *testing.T) {
	v, f := SnowflakeGenerator.Generate()
	if f != nil {
		t.Errorf("%v, %v, %v, %v", f, TimeEpoch, parseerror, snowerror)
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
