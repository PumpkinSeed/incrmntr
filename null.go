package incrmntr

import "time"

type NullTimeout struct {
	valid bool
	Value time.Duration
}

func NullTimeoutMillisec(dur uint64) NullTimeout {
	return NullTimeout{
		valid: true,
		Value: time.Duration(dur) * time.Millisecond,
	}
}

func NullTimeoutSec(dur uint64) NullTimeout {
	return NullTimeout{
		valid: true,
		Value: time.Duration(dur) * time.Second,
	}
}

func NullTimeoutFrom(dur time.Duration) NullTimeout {
	return NullTimeout{
		valid: true,
		Value: dur,
	}
}
