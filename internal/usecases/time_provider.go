package usecases

import "time"

type TimeProvider struct {
}

func (t *TimeProvider) Now() time.Time {
	return time.Now().In(time.Local)
}
