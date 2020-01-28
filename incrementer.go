package incrmntr

import (
	"errors"
	"sync"
	"time"

	"github.com/couchbase/gocb/v2"
)

// Incrmntr is the base interface of the library
type Incrmntr interface {
	Get(key string) (int64, error)
	Add(key string) (NullInt64, error)
	AddSafe(key string) (NullInt64, error)
	AddWithRollover(key string, rollover uint64) (NullInt64, error)
	AddSafeWithRollover(key string, rollover uint64) (NullInt64, error)
	SetTimeout(timeout time.Duration)
	Close() error
}

type BucketOpts struct {
	OperationTimeout      NullTimeout
	BulkOperationTimeout  NullTimeout
	DurabilityTimeout     NullTimeout
	DurabilityPollTimeout NullTimeout
	ViewTimeout           NullTimeout
	N1qlTimeout           NullTimeout
	AnalyticsTimeout      NullTimeout
}

// Incrementer is the main struct stores the related data
// and implements the Incrmntr interface
type Incrementer struct {
	sync.Mutex

	bucket   *gocb.Bucket
	rollover uint64
	initial  int64
	inc      uint64
	cycle    bool
	timeout  time.Duration
}

// New creates a new handler which implements the Incrmntr and setup the buckets
func New(bucket *gocb.Bucket, rollover uint64, initial int64, inc uint64, cycle bool) (Incrmntr, error) {
	return &Incrementer{
		bucket:   bucket,
		rollover: rollover,
		initial:  initial,
		inc:      inc,
		cycle:    cycle,
		timeout:  5000*time.Millisecond,
	}, nil
}

// Get the value of the given key
func (i *Incrementer) Get(key string) (int64, error) {
	var v interface{}
	doc, err := i.bucket.DefaultCollection().Get(key, &gocb.GetOptions{
		Timeout: i.GetTimeout(),
	})
	if err != nil {
		return 0, err
	}
	err = doc.Content(&v)

	return int64(v.(float64)), err
}

// AddWithRollover is do the increment on the specified key
// custom rollover on the key available
func (i *Incrementer) AddWithRollover(key string, rollover uint64) (NullInt64, error) {
	if i.bucket == nil {
		return nullInt64(), errors.New("error bucket is nil")
	}
	return i.add(key, rollover)
}

// AddSafeWithRollover do the increment on the specified key
// concurrency and lock safe increment
// custom rollover on the key available
func (i *Incrementer) AddSafeWithRollover(key string, rollover uint64) (NullInt64, error) {
	if i.bucket == nil {
		return nullInt64(), errors.New("error bucket is nil")
	}

	var value NullInt64
	var err error
	value, err = i.add(key, rollover)
	if errors.Is(err, gocb.ErrTemporaryFailure) {
		for {
			value, err = i.add(key, rollover)
			if err == nil {
				break
			}
		}
	} else if err != nil {
		return nullInt64(), err
	}

	return value, nil
}

// Add is do the increment on the specified key
func (i *Incrementer) Add(key string) (NullInt64, error) {
	if i.bucket == nil {
		return nullInt64(), errors.New("error bucket is nil")
	}
	return i.add(key, i.rollover)
}

// AddSafe do the increment on the specified key
// concurrency and lock safe increment
func (i *Incrementer) AddSafe(key string) (NullInt64, error) {
	if i.bucket == nil {
		return nullInt64(), errors.New("error bucket is nil")
	}

	var value NullInt64
	var err error
	value, err = i.add(key, i.rollover)
	if errors.Is(err, gocb.ErrTemporaryFailure) {
		for {
			value, err = i.add(key, i.rollover)
			if err == nil {
				break
			}
		}
	}else if err != nil {
		return nullInt64(), err
	}

	return value, nil
}

// Close the bucket
func (i *Incrementer) Close() error {
	//err := i.bucket.Close()
	i.bucket = nil
	return nil
}

// add handle the increment mechanism, rollover passed as
// parameter because there is functions with custom rollover
func (i *Incrementer) add(key string, rollover uint64) (NullInt64, error) {
	var err error

	// ---- initKey called first to ensure key will be ready for operation
	initHappened, err := i.initKey(key)
	if err != nil {
		return nullInt64(), err
	}
	if initHappened {
		return nullInt64From(1), nil
	}

	// ---- get the current value and lock the cas
	var current interface{}
	res, err := i.bucket.DefaultCollection().GetAndLock(key, 100*time.Millisecond, &gocb.GetAndLockOptions{
		Timeout: i.GetTimeout(),
	})
	if err != nil {
		return nullInt64(), err
	}
	cas := res.Cas()
	err = res.Content(&current)
	if err != nil {
		return nullInt64(), err
	}

	// ---- do the exact increment mechanism
	newValue := current.(float64) + float64(i.inc)
	if i.cycle && newValue > float64(rollover) {
		newValue = float64(i.initial)
	}

	_, err = i.bucket.DefaultCollection().Replace(key, newValue, &gocb.ReplaceOptions{Expiry: 0, Cas: cas, Timeout: i.GetTimeout()})

	// https://developer.couchbase.com/documentation/server/3.x/developer/dev-guide-3.0/lock-items.html

	return nullInt64From(int64(newValue)), err
}

// initKey do the key initialze process, it's means
// if the key not found, call the Counter which creates it
func (i *Incrementer) initKey(key string) (bool, error) {
	i.Lock()
	defer i.Unlock()

	// ---- is a flag, shows any action happened
	var happened = false

	// ---- check key is exists, if not create it
	_, err := i.bucket.DefaultCollection().Get(key, &gocb.GetOptions{
		Timeout: i.GetTimeout(),
	})
	//res.Content(&v)
	if errors.Is(err, gocb.ErrDocumentNotFound) {
		_, err = i.bucket.DefaultCollection().Binary().Increment(key, &gocb.IncrementOptions{
			Initial: i.initial,
			Delta:   uint64(i.initial),
			Timeout: i.GetTimeout(),
			Expiry:  0, // Seconds
		})
		if err != nil {
			return false, err
		}
		happened = true
	} else {
		return false, err
	}

	return happened, nil
}

func (i *Incrementer) GetTimeout() time.Duration {
	return i.timeout
}

func (i *Incrementer) SetTimeout(timeout time.Duration) {
	i.timeout = timeout
}
