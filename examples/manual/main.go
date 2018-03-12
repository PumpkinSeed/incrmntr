package main

import (
	"fmt"

	"github.com/PumpkinSeed/incrmntr"
)

func main() {
	inc := incrmntr.New("couchbase://localhost", "increment", "", 999999999999999, 1)
	for i := 0; i < 20000; i++ {
		err := inc.AddSafe("test")
		if err != nil {
			fmt.Println(err)
		}
	}
}
