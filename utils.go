package accesslogger

import "time"

var Clock func() time.Time = func() time.Time { return time.Now().Local() }

func coalesce[T comparable](ts ...T) T {
	var empty T
	for _, t := range ts {
		if t != empty {
			return t
		}
	}
	return empty
}

func emptyif[T comparable](t T, v T) T {
	if t != v {
		return t
	}
	var empty T
	return empty
}
