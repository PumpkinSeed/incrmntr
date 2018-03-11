package main

import (
	"fmt"

	"github.com/PumpkinSeed/incrmntr"
	"github.com/couchbase/gocb"
)

func main() {
	inc := incrmntr.New("couchbase://localhost", "increment", "", 999999999999999, 1)
	for i := 0; i < 20000; i++ {
		err := inc.Add("test")
		if err == gocb.ErrTmpFail {
		Loop:
			for {
				err := inc.Add("test")
				if err == nil {
					break Loop
				}
			}
			continue
		}
		if err != gocb.ErrTmpFail && err != nil {
			fmt.Println(err)
			// 3216
		}
	}
}
