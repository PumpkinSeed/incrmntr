package main

import (
	"fmt"

	"github.com/PumpkinSeed/incrmntr"
	"gopkg.in/couchbase/gocb.v1"
)

func main() {
	cluster, err := gocb.Connect("couchbase://localhost")
	if err != nil {
		fmt.Printf("error connecting to the cluster: %s", err.Error())
	}
	cluster.Authenticate(gocb.PasswordAuthenticator{
		Username: "Administrator",
		Password: "password",
	})
	inc, err := incrmntr.New(cluster, "increment", "", 999999999999999, 1)
	if err != nil {
		fmt.Printf("%s", err.Error())
	}
	for i := 0; i < 20000; i++ {
		err := inc.AddSafe("test")
		if err != nil {
			fmt.Println(err)
		}
	}
}
