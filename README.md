# Incrementer

Distributed counter based on Couchbase.

### Usage

Recommendation: Always try to use the `AddSafe`, it's ensure the 100% the key will gets incremented.

[Benchmarks](benchmark.md)

#### Non-persistent connection

- `Add`: Simple add plus one to the specified key, if the key has been locked it will return an `gocb.ErrTmpFail`
- `AddSafe`: Solve the issue occured in the `Add` which cause the `gocb.ErrTmpFail`


**Parameter list of non-persistent functions**

- `conn string`: connection string to couchbase's cluster
- `auth gocb.PasswordAuthenticator`: password auth for cocuhbase's cluster
- `bucketName string`: name of the bucket
- `bucketPassword string`: password of the bucket
- `key string`: key to increment
- `rollover uint64`:  rollover of the key
- `initial int64`: initial value of the key

```
err := incrmntr.Add(
	"couchbase://localhost", 
	gocb.PasswordAuthenticator{
		Username: "Administrator",
		Password: "password",
	}, 
	"increment", 
	"", 
	"test2",
	999,
	1,
)
// handle error
```

#### Persistent connection

**New**

```
New(cluster *gocb.Cluster, bucketName, bucketPassword string, rollover uint64, initial int64) (Incrmntr, error)
```

- `cluster` waits a previously setted up couchbase cluster connection.
- `bucketName` waits the name of the bucket
- `bucketPassword` waits the password of the bucket
- `rollover` is the rollover limit
- `initial` is the initial value if the rollover happens put it back to that number

**Methods**

Interface returned by the `incrmntr.New`

```
type Incrmntr interface {
	Get(key string) (int64, error)
	Add(key string) error
	AddSafe(key string) error
	AddWithRollover(key string, rollover uint64) error
	AddSafeWithRollover(key string, rollover uint64) error
	Close() error
}
```

- `Get`: Returns the current value of the key
- `Add`: Simple add plus one to the specified key, if the key has been locked it will return an `gocb.ErrTmpFail`
- `AddSafe`: Solve the issue occured in the `Add` which cause the `gocb.ErrTmpFail`
- `AddWithRollover`: Same as the `Add` with custom rollover
- `AddSafeWithRollover`: Same as the `AddSafe` with custom rollover
- `Close`: Close the bucket

```
cluster, err := gocb.Connect("couchbase://localhost")
// handle error

cluster.Authenticate(gocb.PasswordAuthenticator{
	Username: "Administrator",
	Password: "password",
})

inc, err := incrmntr.New(cluster, "increment", "", 999, 1)
// handle error

err := inc.Add("test")
// handle error

err := inc.AddSafe("test")
// handle error
```

### Contribution

There is a `docker/docker-compose-single.yml` which represents a single couchbase server

The `docker-compose.yml` has a cluster what the developer have to set up manually. The api server waits a req like that to trigger the test case:

```
curl -X GET \
  'http://localhost:8890/trigger?amount=100&bucket=increment&conn=couchbase://cb2' \
  -H 'content-type: application/x-www-form-urlencoded'
```