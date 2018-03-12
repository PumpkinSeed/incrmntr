package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/PumpkinSeed/incrmntr"
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

	inc := incrmntr.New(conn[0], bucket[0], "", 999999999999999, 1)
	for i := 0; i < amountInt; i++ {
		err := inc.AddSafe("test")
		if err != nil {
			fmt.Println(err)
		}
	}

	fmt.Fprintf(w, "done")
}
