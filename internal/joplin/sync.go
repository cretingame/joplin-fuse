package joplin

import "time"

type Time int

func JoplinTime(t int) time.Time {
	return time.Unix(int64(t/1000), int64((t%1000)*1000_000))
}

func (t Time) JoplinTime() time.Time {
	return JoplinTime(int(t))
}
