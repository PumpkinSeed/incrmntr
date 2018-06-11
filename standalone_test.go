package incrmntr

import (
	"testing"

	"github.com/couchbase/gocb"
)

func TestConAndAdd(t *testing.T) {
	err := Add("couchbase://localhost", gocb.PasswordAuthenticator{
		Username: "Administrator",
		Password: "password",
	}, "increment", "", "test2")

	if err != nil {
		t.Error(err)
	}
}

func BenchmarkConAndAdd(b *testing.B) {

	for i := 0; i < b.N; i++ {
		err := Add("couchbase://localhost", gocb.PasswordAuthenticator{
			Username: "Administrator",
			Password: "password",
		}, "increment", "", "test2")
		if err != nil {
			b.Error(err)
		}
	}
}
