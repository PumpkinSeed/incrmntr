package framework

import (
	"encoding/json"
	"fmt"
	"testing"

	"bitbucket.org/fluidpay/processing-engine/pkg/analytics"

	"github.com/rs/xid"
)

func TestCouchbase(t *testing.T) {
	counter := NewCouchbase(analytics.Default())

	var cfg = config{
		Address:        "couchbase://localhost",
		Username:       "Administrator",
		Password:       "password",
		Bucket:         "increment",
		BucketPassword: "",
		Rollover:       999,
		Initial:        1,
	}
	cfgByte, _ := json.Marshal(cfg)
	fmt.Println(string(cfgByte))

	err := counter.Init(cfgByte)
	if err != nil {
		t.Error(err)
		return
	}

	var iterator = 999
	var val int64
	var key = xid.New().String()
	for i := 0; i < iterator; i++ {
		val, err = counter.NextVal(key)
		if err != nil {
			t.Error(err)
		}
	}

	if val != 999 {
		t.Errorf("Value should be 999, instead of %d", val)
	}
	val, err = counter.NextVal(key)
	if err != nil {
		t.Error(err)
	}
	if val != 1 {
		t.Errorf("Value should be 1, instead of %d", val)
	}
}

/*
	pkg: bitbucket.org/fluidpay/processing-engine/pkg/counter
	BenchmarkCouchbase-4   	     300	   5457130 ns/op

	5.45713 ms/op
*/
func BenchmarkCouchbase(b *testing.B) {
	counter := NewCouchbase(analytics.Default())

	var cfg = config{
		Address:        "couchbase://localhost",
		Username:       "Administrator",
		Password:       "password",
		Bucket:         "increment",
		BucketPassword: "",
		Rollover:       999,
		Initial:        1,
	}
	cfgByte, _ := json.Marshal(cfg)

	err := counter.Init(cfgByte)
	if err != nil {
		b.Error(err)
		return
	}

	var val int64
	var key = xid.New().String()
	for i := 0; i < b.N; i++ {
		val, err = counter.NextVal(key)
		if err != nil {
			b.Error(err)
		}
	}

	b.Log(val)
}
