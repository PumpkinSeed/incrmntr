package incrmntr

import (
	"testing"

	"gopkg.in/couchbase/gocb.v1"
)

func TestConAndAdd(t *testing.T) {
	err := Add("couchbase://localhost", gocb.PasswordAuthenticator{
		Username: "Administrator",
		Password: "password",
	}, "increment", "", "60a5c232-1f1f-40b0-b4f9-cf51808d96eb-conandadd", 999, 1, 1)

	if err != nil {
		t.Error(err)
	}
}

func TestConAndAddSafe(t *testing.T) {
	err := AddSafe("couchbase://localhost", gocb.PasswordAuthenticator{
		Username: "Administrator",
		Password: "password",
	}, "increment", "", "60a5c232-1f1f-40b0-b4f9-cf51808d96eb-conandadd", 999, 1, 1)

	if err != nil {
		t.Error(err)
	}
}

func BenchmarkConAndAdd(b *testing.B) {
	for i := 0; i < b.N; i++ {
		err := Add("couchbase://localhost", gocb.PasswordAuthenticator{
			Username: "Administrator",
			Password: "password",
		}, "increment", "", "3cf323ec-79b4-4b67-9f52-2655d1227e71-conandadd", 999, 1, 1)
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkConAndAddSafe(b *testing.B) {
	for i := 0; i < b.N; i++ {
		err := AddSafe("couchbase://localhost", gocb.PasswordAuthenticator{
			Username: "Administrator",
			Password: "password",
		}, "increment", "", "3cf323ec-79b4-4b67-9f52-2655d1227e71-conandadd", 999, 1, 1)
		if err != nil {
			b.Error(err)
		}
	}
}
