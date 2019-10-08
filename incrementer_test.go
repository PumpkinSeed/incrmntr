package incrmntr

import (
	"github.com/rs/xid"
	"sync"
	"testing"

	"github.com/couchbase/gocb"
)

var skipTest = map[string]bool{
	"add":                 false,
	"addsafe":             false,
	"addwithrollover":     false,
	"addsafewithrollover": false,
	"initkey":             false,
}

func TestAdd(t *testing.T) {
	if skipTest["add"] {
		t.Skip("Add skipped")
	}

	var rollover = int64(999)
	var init = int64(1)
	var key = xid.New().String()
	var testCounter = newCounterTest(init, rollover)

	cluster, err := gocb.Connect("couchbase://localhost")
	if err != nil {
		t.Errorf("error connecting to the cluster: %s", err.Error())
	}
	cluster.Authenticate(gocb.PasswordAuthenticator{
		Username: "Administrator",
		Password: "password",
	})
	// Open Bucket
	bucket, err := cluster.OpenBucket("increment", "")
	if err != nil {
		t.Error(err)
	}

	i, err := New(bucket, uint64(rollover), init, 1)
	if err != nil {
		t.Error(err)
	}

	i.Add(key)
	testCounter.add()
	i.Add(key)
	testCounter.add()
	i.Add(key)
	testCounter.add()
	val, err := i.Get(key)
	if err != nil {
		t.Error(err)
	}
	if val != testCounter.val {
		t.Errorf("Incrementer value should be %d, instead of %d", testCounter.val, val)
	}
}

func TestAddSafe(t *testing.T) {
	if skipTest["addsafe"] {
		t.Skip("AddSafe skipped")
	}

	var rollover = int64(99)
	var init = int64(1)
	var key = xid.New().String()
	var testCounter = newCounterTest(init, rollover)

	cluster, err := gocb.Connect("couchbase://localhost")
	if err != nil {
		t.Errorf("error connecting to the cluster: %s", err.Error())
	}
	cluster.Authenticate(gocb.PasswordAuthenticator{
		Username: "Administrator",
		Password: "password",
	})
	// Open Bucket
	bucket, err := cluster.OpenBucket("increment", "")
	if err != nil {
		t.Error(err)
	}

	inc, err := New(bucket, uint64(rollover), init, 1)
	if err != nil {
		t.Error(err)
	}
	var wg sync.WaitGroup

	for i := 0; i < 103; i++ {
		wg.Add(1)
		go func() {
			_, err := inc.AddSafe(key)
			testCounter.add()
			if err != nil {
				t.Error(err)
			}
			wg.Done()
		}()
	}
	wg.Wait()
	val, err := inc.Get(key)
	if err != nil {
		t.Error(err)
	}
	if val != testCounter.val {
		t.Errorf("Incrementer value should be %d, instead of %d", testCounter.val, val)
	}
}

func TestAddWithRollover(t *testing.T) {
	if skipTest["addwithrollover"] {
		t.Skip("AddWithRollover skipped")
	}

	var rollover = int64(99)
	var init = int64(1)
	var key = xid.New().String()
	var testCounter = newCounterTest(init, rollover)

	cluster, err := gocb.Connect("couchbase://localhost")
	if err != nil {
		t.Errorf("error connecting to the cluster: %s", err.Error())
	}
	cluster.Authenticate(gocb.PasswordAuthenticator{
		Username: "Administrator",
		Password: "password",
	})

	// Open Bucket
	bucket, err := cluster.OpenBucket("increment", "")
	if err != nil {
		t.Error(err)
	}

	i, err := New(bucket, uint64(rollover), init, 1)
	if err != nil {
		t.Error(err)
	}

	i.AddWithRollover(key, 23)
	testCounter.addWithRollover(23)
	i.AddWithRollover(key, 23)
	testCounter.addWithRollover(23)
	i.AddWithRollover(key, 23)
	testCounter.addWithRollover(23)
	val, err := i.Get(key)
	if err != nil {
		t.Error(err)
	}
	if val != testCounter.val {
		t.Errorf("Incrementer value should be %d, instead of %d", testCounter.val, val)
	}
}

func TestAddSafeWithRollover(t *testing.T) {
	if skipTest["addsafewithrollover"] {
		t.Skip("AddSafeWithRollover skipped")
	}

	var rollover = int64(99)
	var init = int64(1)
	var key = xid.New().String()
	var testCounter = newCounterTest(init, rollover)

	cluster, err := gocb.Connect("couchbase://localhost")
	if err != nil {
		t.Errorf("error connecting to the cluster: %s", err.Error())
	}
	cluster.Authenticate(gocb.PasswordAuthenticator{
		Username: "Administrator",
		Password: "password",
	})
	// Open Bucket
	bucket, err := cluster.OpenBucket("increment", "")
	if err != nil {
		t.Error(err)
	}

	inc, err := New(bucket, uint64(rollover), init, 1)
	if err != nil {
		t.Error(err)
	}
	var wg sync.WaitGroup

	for i := 0; i < 103; i++ {
		wg.Add(1)
		go func() {
			_, err := inc.AddSafeWithRollover(key, 55)
			if err != nil {
				t.Error(err)
			}
			testCounter.addWithRollover(55)
			wg.Done()
		}()
	}
	wg.Wait()
	val, err := inc.Get(key)
	if err != nil {
		t.Error(err)
	}
	if val != testCounter.val {
		t.Errorf("Incrementer value should be %d, instead of %d", testCounter.val, val)
	}
}

func TestIncrementer_ReturnAdd(t *testing.T) {
	var rollover = int64(99)
	var init = int64(1)
	var key = xid.New().String()

	cluster, err := gocb.Connect("couchbase://localhost")
	if err != nil {
		t.Errorf("error connecting to the cluster: %s", err.Error())
	}
	cluster.Authenticate(gocb.PasswordAuthenticator{
		Username: "Administrator",
		Password: "password",
	})
	// Open Bucket
	bucket, err := cluster.OpenBucket("increment", "")
	if err != nil {
		t.Error(err)
	}

	inc, err := New(bucket, uint64(rollover), init, 1)
	if err != nil {
		t.Error(err)
	}

	var value NullInt64
	for i := 0; i<10;i++ {
		var err error
		value, err = inc.AddSafe(key)
		if err != nil {
			t.Fatal(err)
		}
	}

	if value.Valid {
		if value.Value != 10 {
			t.Errorf("Value should be 10, instead of %d", value.Value)
		}
	} else {
		t.Error("Value should be valid")
	}
}

func TestInitKey(t *testing.T) {
	if skipTest["initkey"] {
		t.Skip("InitKey skipped")
	}
	var key = xid.New().String()

	cluster, err := gocb.Connect("couchbase://localhost")
	if err != nil {
		t.Errorf("error connecting to the cluster: %s", err.Error())
	}
	cluster.Authenticate(gocb.PasswordAuthenticator{
		Username: "Administrator",
		Password: "password",
	})
	// Open Bucket
	bucket, err := cluster.OpenBucket("increment", "")
	if err != nil {
		t.Error(err)
	}

	inc, err := New(bucket, 99, 1, 1)
	if err != nil {
		t.Error(err)
	}

	incrementer := inc.(*Incrementer)
	_, err = incrementer.initKey(key)
	if err != nil {
		t.Error(err)
	}
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
	// Open Bucket
	bucket, err := cluster.OpenBucket("increment", "")
	if err != nil {
		b.Error(err)
	}

	inc, err := New(bucket, 999, 1, 1)
	if err != nil {
		b.Error(err)
	}

	for i := 0; i < b.N; i++ {
		_, err := inc.Add("b88c972c-e7a8-4d47-a67a-5c7f89914595-b-add")
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkAddSafe(b *testing.B) {
	cluster, err := gocb.Connect("couchbase://localhost")
	if err != nil {
		b.Errorf("error connecting to the cluster: %s", err.Error())
	}
	cluster.Authenticate(gocb.PasswordAuthenticator{
		Username: "Administrator",
		Password: "password",
	})

	// Open Bucket
	bucket, err := cluster.OpenBucket("increment", "")
	if err != nil {
		b.Error(err)
	}

	inc, err := New(bucket, 999, 1, 1)
	if err != nil {
		b.Error(err)
	}

	for i := 0; i < b.N; i++ {
		_, err := inc.AddSafe("b88c972c-e7a8-4d47-a67a-5c7f89914595-b-addsafe")
		if err != nil {
			b.Error(err)
		}
	}
}

/*
	represent real value
*/

type counterTest struct {
	sync.RWMutex
	init     int64
	rollover int64
	val      int64
}

func newCounterTest(init int64, rollover int64) counterTest {
	return counterTest{
		init:     init,
		rollover: rollover,
		val:      init - 1,
	}
}

func (c *counterTest) add() {
	c.Lock()
	defer c.Unlock()
	c.val++
	if c.val > c.rollover {
		c.val = c.init
		return
	}

	return
}

func (c *counterTest) addWithRollover(rollover int64) {
	c.Lock()
	defer c.Unlock()
	c.val++
	if c.val > rollover {
		c.val = c.init
		return
	}

	return
}
