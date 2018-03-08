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

func New(conn, bucketName, bucketPassword string, gap uint64, initial int64) Incrementer {
	cluster, err := gocb.Connect("couchbase://127.0.0.1")
	if err != nil {
		fmt.Println("ERRROR CONNECTING TO CLUSTER:", err)
	}

	// Open Bucket
	bucket, err := cluster.OpenBucket(bucketName, bucketPassword)
	if err != nil {
		fmt.Println("ERRROR OPENING BUCKET:", err)
	}

	return Incrementer{
		bucket:  bucket,
		gap:     gap,
		initial: initial,
	}
}

func (i *Incrementer) Add(key string) error {
	var current interface{}
	cas, err := i.bucket.GetAndLock(key, i.ttl, &current)
	if err == gocb.ErrKeyNotFound {
		i.bucket.Counter(key, 1, i.initial, 0)
	}
	if err != nil {
		return err
	}
	newValue := current.(float64) + 1
	if newValue >= float64(i.gap) {
		newValue = float64(i.initial)
	}
	_, err = i.bucket.Replace(key, newValue, cas, 0)

	/*current, _, err := i.bucket.Counter(key, 1, i.initial, 0)
	if current >= i.gap+1 {
		_, err = i.bucket.Upsert(key, i.initial, 0)
		return err
	}*/
	//fmt.Printf("Current value: %d\n", curKeyValue)
	return err
}
