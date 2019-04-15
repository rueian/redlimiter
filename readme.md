# redlimiter

[![CircleCI](https://circleci.com/gh/rueian/redlimiter.svg?style=svg)](https://circleci.com/gh/rueian/redlimiter)
[![codecov](https://codecov.io/gh/rueian/redlimiter/branch/master/graph/badge.svg)](https://codecov.io/gh/rueian/redlimiter)

A golang library based on go-redis connection pool that can be used to implement rate limit function.

## Counter

The `counter` package provides 2 `Incrementer` implementations
that both `INCR` a key in redis and set expiration atomically when the key created without race condition.

```go
package main

import (
	"github.com/go-redis/redis"
	"github.com/rueian/redlimiter/counter"
)

func main() {
	client := redis.NewUniversalClient(&redis.Options{Addr: "localhost:6379"})
	
	incrementer := counter.NewLuaIncrementer(client) // or counter.NewTxIncrementer(client)
	count, err := incrementer.Incr("client_ip", 1) // expire in 1 sec
	if err != nil {
		rejectWithErr(err)
	}
	if count > 50 {
		rejectWithLimitExceeded()
	}
}
```