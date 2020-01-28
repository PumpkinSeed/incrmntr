package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/PumpkinSeed/incrmntr"
	"github.com/couchbase/gocb/v2"
	"github.com/pkg/profile"
)

const iterationAdd = 1000
const iterationConAndAdd = 50

func main() {
	var standalone bool
	var memProf bool
	flag.BoolVar(&standalone, "standalone", false, "Standalone mode")
	flag.BoolVar(&memProf, "mem-profile", false, "Mem profile")
	flag.Parse()

	if memProf {
		defer profile.Start(profile.MemProfile, profile.ProfilePath("."), profile.NoShutdownHook).Stop()
	} else {
		defer profile.Start(profile.CPUProfile, profile.ProfilePath("."), profile.NoShutdownHook).Stop()
	}

	if standalone {
		//fmt.Println(standaloneAdd())
		return
	}

	fmt.Println(add())
	return
}

//func standaloneAdd() error {
//	for i := 0; i < iterationConAndAdd; i++ {
//		err := incrmntr.Add(
//			"couchbase://localhost",
//			gocb.PasswordAuthenticator{
//				Username: "Administrator",
//				Password: "password",
//			},
//			"increment",
//			"",
//			"test2",
//		)
//		if err != nil {
//			return err
//		}
//	}
//	return nil
//}

func add() error {
	opts := gocb.ClusterOptions{
		TimeoutsConfig: gocb.TimeoutsConfig{KVTimeout: 10*time.Second, QueryTimeout: 10*time.Second},
		Authenticator: gocb.PasswordAuthenticator{
			"Administrator",
			"password",
		},
	}
	cluster, err := gocb.Connect("localhost", opts)
	if err != nil {
		panic(err)
	}

	// get a bucket reference
	bucket := cluster.Bucket("increment")

	//inc := New("couchbase://cb1,cb2", "increment", "", 999, 1)
	inc, err := incrmntr.New(bucket, 999,1 , 1, true)
	if err != nil {
		return err
	}

	for i := 0; i < iterationAdd; i++ {
		_, err := inc.Add("test")
		if err != nil {
			return err
		}
	}

	return nil
}
