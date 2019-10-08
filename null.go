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

type NullInt64 struct {
	Valid bool
	Value int64
}

func nullInt64() NullInt64 {
	return NullInt64{
		Valid: false,
	}
}

func nullInt64From(v int64) NullInt64 {
	return NullInt64{
		Valid: true,
		Value: v,
	}
}
