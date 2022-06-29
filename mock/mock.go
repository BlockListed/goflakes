package mock

import "time"

type TimeNow func() time.Time

var MockTimeNow TimeNow

func CurryMockTimeNow(t time.Time) TimeNow {
	return func() time.Time {
		return t
	}
}

func Time_Now() time.Time {
	if MockTimeNow != nil {
		return MockTimeNow()
	} else {
		return time.Now()
	}
}
