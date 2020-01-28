package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/couchbase/gocb/v2"
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

	//conn, ok := r.Form["conn"]
	//if !ok {
	//	fmt.Fprintf(w, "conn not provided")
	//	return
	//}
	//
	//bucket, ok := r.Form["bucket"]
	//if !ok {
	//	fmt.Fprintf(w, "bucket not provided")
	//	return
	//}

	//amount, ok := r.Form["amount"]
	//if !ok {
	//	fmt.Fprintf(w, "amount not provided")
	//	return
	//}
	//
	//amountInt, err := strconv.Atoi(amount[0])
	//if err != nil {
	//	fmt.Fprintf(w, err.Error())
	//	return
	//}
	//
	//bucketConn, closeC := getBucket()
	//defer closeC(nil)

	//inc, err := incrmntr.New(bucketConn, 999, 1, 1)
	//if err != nil {
	//	fmt.Fprintf(w, err.Error())
	//	return
	//}
	//for i := 0; i < amountInt; i++ {
	//	err := inc.AddSafe("test")
	//	if err != nil {
	//		fmt.Fprintf(w, err.Error())
	//		return
	//	}
	//}

	fmt.Fprintf(w, "done")
}

func getBucket() (*gocb.Bucket, func(opts *gocb.ClusterCloseOptions) error) {
	opts := gocb.ClusterOptions{
		TimeoutsConfig: gocb.TimeoutsConfig{KVTimeout: 10 * time.Second, QueryTimeout: 10 * time.Second},
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
	return cluster.Bucket("increment"), cluster.Close
}
