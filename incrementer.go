package incrmntr

import (
	"errors"
	"fmt"
	"sync"

	"github.com/couchbase/gocb"
)

// Incrmntr is the base interface of the library
type Incrmntr interface {
	Get(key string) (int64, error)
	Add(key string) error
	AddSafe(key string) error
	AddWithRollover(key string, rollover uint64) error
	AddSafeWithRollover(key string, rollover uint64) error
	Close() error
}

// Incrementer is the main struct stores the related data
// and implements the Incrmntr interface
type Incrementer struct {
	sync.Mutex

	bucket   *gocb.Bucket
	rollover uint64
	initial  int64
	ttl      uint32
}

// New creates a new handler which implements the Incrmntr and setup the buckets
func New(cluster *gocb.Cluster, bucketName, bucketPassword string, rollover uint64, initial int64) (Incrmntr, error) {
	// Open Bucket
	bucket, err := cluster.OpenBucket(bucketName, bucketPassword)
	if err != nil {
		return nil, fmt.Errorf("error opening the bucket: %s", err.Error())
	}

	return &Incrementer{
		bucket:   bucket,
		rollover: rollover,
		initial:  initial,
	}, nil
}

// Get the value of the given key
func (i *Incrementer) Get(key string) (int64, error) {
	var v interface{}
	_, err := i.bucket.Get(key, &v)

	return int64(v.(float64)), err
}

// AddWithRollover is do the increment on the specified key
// custom rollover on the key available
func (i *Incrementer) AddWithRollover(key string, rollover uint64) error {
	if i.bucket == nil {
		return errors.New("error bucket is nil")
	}
	return i.add(key, rollover)
}

// AddSafeWithRollover do the increment on the specified key
// concurrency and lock safe increment
// custom rollover on the key available
func (i *Incrementer) AddSafeWithRollover(key string, rollover uint64) error {
	if i.bucket == nil {
		return errors.New("error bucket is nil")
	}

	err := i.add(key, rollover)
	if err == gocb.ErrTmpFail {
		for {
			err := i.add(key, rollover)
			if err == nil {
				break
			}
		}
	}
	if err != gocb.ErrTmpFail && err != nil {
		return err
	}

	return nil
}

// Add is do the increment on the specified key
func (i *Incrementer) Add(key string) error {
	if i.bucket == nil {
		return errors.New("error bucket is nil")
	}
	return i.add(key, i.rollover)
}

// AddSafe do the increment on the specified key
// concurrency and lock safe increment
func (i *Incrementer) AddSafe(key string) error {
	if i.bucket == nil {
		return errors.New("error bucket is nil")
	}

	err := i.add(key, i.rollover)
	if err == gocb.ErrTmpFail {
		for {
			err := i.add(key, i.rollover)
			if err == nil {
				break
			}
		}
	}
	if err != gocb.ErrTmpFail && err != nil {
		return err
	}

	return nil
}

// Close the bucket
func (i *Incrementer) Close() error {
	err := i.bucket.Close()
	i.bucket = nil
	return err
}

// add handle the increment mechanism, rollover passed as
// parameter because there is functions with custom rollover
func (i *Incrementer) add(key string, rollover uint64) error {
	var err error

	// ---- initKey called first to ensure key will be ready for operation
	initHappened, err := i.initKey(key)
	if err != nil {
		return err
	}
	if initHappened {
		return nil
	}

	// ---- get the current value and lock the cas
	var current interface{}
	cas, err := i.bucket.GetAndLock(key, i.ttl, &current)
	if err != nil {
		return err
	}

	// ---- do the exact increment mechanism
	newValue := current.(float64) + 1
	if newValue > float64(rollover) {
		newValue = float64(i.initial)
	}
	_, err = i.bucket.Replace(key, newValue, cas, 0)

	// https://developer.couchbase.com/documentation/server/3.x/developer/dev-guide-3.0/lock-items.html

	return err
}

// initKey do the key initialze process, it's means
// if the key not found, call the Counter which creates it
func (i *Incrementer) initKey(key string) (bool, error) {
	i.Lock()
	defer i.Unlock()

	// ---- v stores the value of the key
	var v interface{}

	// ---- is a flag, shows any action happened
	var happened = false

	// ---- check key is exists, if not create it
	_, err := i.bucket.Get(key, &v)
	if err == gocb.ErrKeyNotFound {
		i.bucket.Counter(key, i.initial, i.initial, 0)
		happened = true
	} else {
		return false, err
	}

	// ---- if action happened then check it's valid and return nil
	if happened {
		_, err = i.bucket.Get(key, &v)
		if err == nil && int64(v.(float64)) == i.initial {
			return true, nil
		}
		return false, err
	}

	return false, nil
}
