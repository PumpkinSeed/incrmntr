package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/PumpkinSeed/incrmntr/framework"
)

var load = 20

func main() {
	counterConfig := []byte(`{"address":"couchbase://localhost","username":"Administrator","password":"asd12345","bucket":"increment","bucket_password":"","rollover":999,"initial":1}`)
	c := framework.NewCouchbase()
	err := c.Init(counterConfig)
	if err != nil {
		log.Fatal(err)
	}

	var val int64
	var wg = new(sync.WaitGroup)
	for i := 0; i < load; i++ {
		wg.Add(1)
		go func() {
			t := time.Now()
			val, err = c.NextVal("test")
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(time.Since(t))
			wg.Done()
		}()

	}
	wg.Wait()
	fmt.Println(val)
}
