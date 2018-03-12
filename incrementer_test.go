package incrmntr

import "testing"

func TestAdd(t *testing.T) {
	//i := New("couchbase://cb1,cb2", "increment", "", 999, 1)
	i, err := New("couchbase://localhost", "increment", "", 999, 1)
	if err != nil {
		t.Error(err)
	}
	i.Add("test")
	i.Add("test")
	i.Add("test")
}

func BenchmarkAdd(b *testing.B) {

	//inc := New("couchbase://cb1,cb2", "increment", "", 999, 1)
	inc, err := New("couchbase://localhost", "increment", "", 999, 1)
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
