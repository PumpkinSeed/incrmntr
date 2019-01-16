package framework

import (
	"encoding/json"
	"errors"

	"github.com/PumpkinSeed/incrmntr"
	"github.com/couchbase/gocb"
)

// Counter is the main definition of the counter
type Counter interface {
	Init(config []byte) error
	NextVal(key string) (int64, error)
	NextValWithRollover(key string, rollover uint64) (int64, error)
	Stop() error
}

// config is the main config of the couchbase
type config struct {
	Address        string `json:"address"`
	Username       string `json:"username"`
	Password       string `json:"password"`
	Bucket         string `json:"bucket"`
	BucketPassword string `json:"bucket_password"`
	Rollover       uint64 `json:"rollover"`
	Initial        int64  `json:"initial"`
}

// couchbase is the implementation of Counter with couchbase
type couchbase struct {
	inc incrmntr.Incrmntr

	// mut *sync.Mutex
}

// NewCouchbase creates a new implementation of Counter with couchbase
func NewCouchbase() Counter {
	return &couchbase{
		//analytics: analytics,
		// mut: &sync.Mutex{},
	}
}

// Init an incrementer based on the config
func (c *couchbase) Init(cfgByte []byte) error {
	var cfg config
	err := json.Unmarshal(cfgByte, &cfg)
	if err != nil {
		return err
	}

	cluster, err := gocb.Connect(cfg.Address)
	if err != nil {
		return err
	}

	cluster.Authenticate(gocb.PasswordAuthenticator{
		Username: cfg.Username,
		Password: cfg.Password,
	})

	c.inc, err = incrmntr.New(cluster, cfg.Bucket, cfg.BucketPassword, cfg.Rollover, cfg.Initial)
	if err != nil {
		return err
	}

	return nil
}

// NextVal returns the next value of the key
func (c *couchbase) NextVal(key string) (int64, error) {
	if c.inc == nil {
		return 0, errors.New("nil increment")
	}
	// c.mut.Lock()
	// defer c.mut.Unlock()
	err := c.inc.AddSafe(key)
	if err != nil {
		return 0, err
	}

	return c.inc.Get(key)
}

// NextValWithRollover returns the next value of the key with rollover
func (c *couchbase) NextValWithRollover(key string, rollover uint64) (int64, error) {
	if c.inc == nil {
		return 0, errors.New("nil increment")
	}
	// c.mut.Lock()
	// defer c.mut.Unlock()
	err := c.inc.AddSafeWithRollover(key, rollover)
	if err != nil {
		return 0, err
	}

	return c.inc.Get(key)
}

func (c *couchbase) Stop() error {
	return c.inc.Close()
}
