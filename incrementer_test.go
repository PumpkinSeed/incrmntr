package incrmntr

import (
	"testing"

	"github.com/couchbase/gocb"
)

func TestAdd(t *testing.T) {
	cluster, err := gocb.Connect("couchbase://localhost")
	if err != nil {
		t.Errorf("error connecting to the cluster: %s", err.Error())
	}
	cluster.Authenticate(gocb.PasswordAuthenticator{
		Username: "Administrator",
		Password: "password",
	})
	//i := New("couchbase://cb1,cb2", "increment", "", 999, 1)
	i, err := New(cluster, "increment", "", 999, 1)
	if err != nil {
		t.Error(err)
	}
	i.Add("test")
	i.Add("test")
	i.Add("test")
}

func BenchmarkAdd(b *testing.B) {
	cluster, err := gocb.Connect("couchbase://localhost")
	if err != nil {
		b.Errorf("error connecting to the cluster: %s", err.Error())
	}
	cluster.Authenticate(gocb.PasswordAuthenticator{
		Username: "Administrator",
		Password: "password",
	})
	//inc := New("couchbase://cb1,cb2", "increment", "", 999, 1)
	inc, err := New(cluster, "increment", "", 999, 1)
	if err != nil {
		b.Error(err)
	}

	for i := 0; i < b.N; i++ {
		err := inc.Add("test")
		if err != nil {
			b.Error(err)
		}
	}
}
