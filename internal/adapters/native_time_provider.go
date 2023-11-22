package adapters

import "time"

type NativeTimeProvider struct {
}

func (t *NativeTimeProvider) Now() time.Time {
	return time.Now().In(time.Local)
}
