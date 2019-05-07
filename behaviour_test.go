package incrmntr

import (
	"sync"
	"testing"

	"gopkg.in/couchbase/gocb.v1"
)

func TestBehaviour(t *testing.T) {
	var wg sync.WaitGroup

	wg.Add(1)
	go flow1(t, &wg)
	wg.Add(1)
	go flow2(t, &wg)

	wg.Wait()
}

func flow1(t *testing.T, wg *sync.WaitGroup) {
	var rollover = uint64(999)
	var init = int64(1)
	var keys = []string{
		"9559d644-3a8e-4589-99af-85f8be3e49a3-flow1",
		"e0d161e4-9b8b-431b-8177-dca40b4c80a7-flow1",
		"afb7f56a-6866-4d00-a1c4-031f909f1a18-flow1",
		"2db2e29f-d790-475b-b727-6d89fe38b8d2-common",
	}

	cluster, err := gocb.Connect("couchbase://localhost")
	if err != nil {
		t.Errorf("error connecting to the cluster: %s", err.Error())
	}
	cluster.Authenticate(gocb.PasswordAuthenticator{
		Username: "Administrator",
		Password: "password",
	})
	inc, err := New(cluster, "increment", "", rollover, init, 1)
	if err != nil {
		t.Error(err)
	}

	var wgInner sync.WaitGroup

	for i := 0; i < 103; i++ {
		wgInner.Add(1)
		go func() {
			for _, key := range keys {
				err := inc.AddSafe(key)
				if err != nil {
					t.Error(err)
				}
			}
			wgInner.Done()
		}()
	}
	wgInner.Wait()

	wg.Done()
}

func flow2(t *testing.T, wg *sync.WaitGroup) {
	var rollover = uint64(999)
	var init = int64(1)
	var keys = []string{
		"65ffc8a4-4ea2-47b0-b87e-9a08cd1bf4b7-flow2",
		"04ae092e-6ada-4b3f-a0e6-36f215ad38c7-flow2",
		"b2c61051-e387-47d7-b1f4-1dada71495b8-flow2",
		"2db2e29f-d790-475b-b727-6d89fe38b8d2-common",
	}

	cluster, err := gocb.Connect("couchbase://localhost")
	if err != nil {
		t.Errorf("error connecting to the cluster: %s", err.Error())
	}
	cluster.Authenticate(gocb.PasswordAuthenticator{
		Username: "Administrator",
		Password: "password",
	})
	inc, err := New(cluster, "increment", "", rollover, init, 1)
	if err != nil {
		t.Error(err)
	}

	var wgInner sync.WaitGroup

	for i := 0; i < 103; i++ {
		wgInner.Add(1)
		go func() {
			for _, key := range keys {
				err := inc.AddSafe(key)
				if err != nil {
					t.Error(err)
				}
			}
			wgInner.Done()
		}()
	}
	wgInner.Wait()

	wg.Done()
}
