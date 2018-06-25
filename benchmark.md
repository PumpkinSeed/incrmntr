```
goos: darwin
goarch: amd64
pkg: github.com/PumpkinSeed/incrmntr
BenchmarkAdd-4             	    1000	   2365174 ns/op
BenchmarkAddSafe-4         	    1000	   3720169 ns/op
BenchmarkConAndAdd-4       	      30	  55185922 ns/op
BenchmarkConAndAddSafe-4   	      30	  50782230 ns/op
PASS
ok  	github.com/PumpkinSeed/incrmntr	42.845s
```