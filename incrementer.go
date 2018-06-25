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

func (i *Incrementer) Get(key string) (int64, error) {
	i.Lock()
	defer i.Unlock()

	var v interface{}
	_, err := i.bucket.Get(key, &v)

	return int64(v.(float64)), err
}

func (i *Incrementer) AddWithRollover(key string, rollover uint64) error {
	if i.bucket == nil {
		return errors.New("error bucket is nil")
	}
	return i.add(key, rollover)
}

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
	return i.bucket.Close()
}

type initHappened struct {
	val      uint64
	happened bool
}

func (i *Incrementer) add(key string, rollover uint64) error {
	var err error
	err = i.initKey(key)
	if err != nil {
		return err
	}

	var current interface{}
	cas, err := i.bucket.GetAndLock(key, i.ttl, &current)
	if err != nil {
		return err
	}

	i.Lock()
	defer i.Unlock()
	newValue := current.(float64) + 1
	if newValue > float64(rollover) {
		newValue = float64(i.initial)
	}
	_, err = i.bucket.Replace(key, newValue, cas, 0)

	// https://developer.couchbase.com/documentation/server/3.x/developer/dev-guide-3.0/lock-items.html

	return err
}

func (i *Incrementer) initKey(key string) error {
	i.Lock()
	defer i.Unlock()

	var v interface{}
	var happened = false

	_, err := i.bucket.Get(key, &v)
	if err == gocb.ErrKeyNotFound {
		i.bucket.Counter(key, i.initial, i.initial, 0)
		happened = true
	} else {
		return err
	}

	if happened {
		_, err = i.bucket.Get(key, &v)
		if err != nil && v.(int64) == i.initial {
			return nil
		}
		return err
	}

	return nil
}
