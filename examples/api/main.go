package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/PumpkinSeed/incrmntr"
	"github.com/couchbase/gocb"
)

func main() {
	http.HandleFunc("/trigger", trigger)     // set router
	err := http.ListenAndServe(":9090", nil) // set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func trigger(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	conn, ok := r.Form["conn"]
	if !ok {
		fmt.Fprintf(w, "conn not provided")
		return
	}

	bucket, ok := r.Form["bucket"]
	if !ok {
		fmt.Fprintf(w, "bucket not provided")
		return
	}

	amount, ok := r.Form["amount"]
	if !ok {
		fmt.Fprintf(w, "amount not provided")
		return
	}

	amountInt, err := strconv.Atoi(amount[0])
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	cluster, err := gocb.Connect(conn[0])
	if err != nil {
		fmt.Fprintf(w, "error connecting to the cluster: %s", err.Error())
		return
	}
	cluster.Authenticate(gocb.PasswordAuthenticator{
		Username: "Administrator",
		Password: "password",
	})
	bucketConn, err := cluster.OpenBucket(bucket[0], "")
	if err != nil {
		panic(err)
	}

	inc, err := incrmntr.New(bucketConn, 999, 1, 1)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	for i := 0; i < amountInt; i++ {
		err := inc.AddSafe("test")
		if err != nil {
			fmt.Fprintf(w, err.Error())
			return
		}
	}

	fmt.Fprintf(w, "done")
}
