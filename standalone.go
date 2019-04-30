package incrmntr

import (
	"fmt"

	"github.com/couchbase/gocb"
)

// Add implements a standalone version off the incrementer's
// Add method with connect to the Cocuhabse and Close cluster and bucket
func Add(conn string, auth gocb.PasswordAuthenticator, bucketName, bucketPassword string, key string, rollover uint64, initial int64) error {
	cluster, err := gocb.Connect(conn)
	if err != nil {
		return fmt.Errorf("error connecting to the cluster: %s", err.Error())
	}
	//defer cluster.Close()

	cluster.Authenticate(auth)
	i, err := New(cluster, "increment", "", rollover, initial)
	if err != nil {
		return err
	}

	err = i.Add(key)
	if err != nil {
		return err
	}

	return i.Close()
}

// AddSafe implements a standalone version off the incrementer's
// AddSafe method with connect to the Cocuhabse and Close cluster and bucket
func AddSafe(conn string, auth gocb.PasswordAuthenticator, bucketName, bucketPassword string, key string, rollover uint64, initial int64) error {
	cluster, err := gocb.Connect(conn)
	if err != nil {
		return fmt.Errorf("error connecting to the cluster: %s", err.Error())
	}
	//defer cluster.Close()

	cluster.Authenticate(auth)
	i, err := New(cluster, "increment", "", rollover, initial)
	if err != nil {
		return err
	}

	err = i.AddSafe(key)
	if err != nil {
		return err
	}

	return i.Close()
}
