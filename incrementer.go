package incrmntr

import (
	"fmt"

	"github.com/couchbase/gocb"
)

type Incrementer struct {
	bucket  *gocb.Bucket
	gap     uint64
	initial int64
	ttl     uint32
}

func New(conn, bucketName, bucketPassword string, gap uint64, initial int64) *Incrementer {
	cluster, err := gocb.Connect(conn)
	if err != nil {
		fmt.Println("ERRROR CONNECTING TO CLUSTER:", err)
		return nil
	}

	// Open Bucket
	bucket, err := cluster.OpenBucket(bucketName, bucketPassword)
	if err != nil {
		fmt.Println("ERRROR OPENING BUCKET:", err)
		return nil
	}

	return &Incrementer{
		bucket:  bucket,
		gap:     gap,
		initial: initial,
	}
}

func (i *Incrementer) Add(key string) error {
	return i.add(key)
}

func (i *Incrementer) AddSafe(key string) error {
	err := i.add(key)
	if err == gocb.ErrTmpFail {
		for {
			err := i.add("key")
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

func (i *Incrementer) add(key string) error {
	var current interface{}
	cas, err := i.bucket.GetAndLock(key, i.ttl, &current)
	if err == gocb.ErrKeyNotFound {
		_, _, err := i.bucket.Counter(key, 1, i.initial, 0)
		if err != nil {
			return err
		}
	}
	if err != nil {
		return err
	}
	newValue := current.(float64) + 1
	if newValue >= float64(i.gap) {
		newValue = float64(i.initial)
	}
	_, err = i.bucket.Replace(key, newValue, cas, 0)

	// https://developer.couchbase.com/documentation/server/3.x/developer/dev-guide-3.0/lock-items.html

	return err
}
