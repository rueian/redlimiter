package counter

import (
	"sync"
	"testing"
	"time"

	"github.com/go-redis/redis"
)

func connect() *redis.Client {
	return redis.NewClient(&redis.Options{Addr: "redis:6379"})
}

func TestLuaIncrementer(t *testing.T) {
	testIncrementer(t, "lua", NewLuaIncrementer(connect()))
}

func TestTxIncrementer(t *testing.T) {
	testIncrementer(t, "tx", NewTxIncrementer(connect()))
}

func BenchmarkLuaIncrementer(b *testing.B) {
	benchIncrementer(b, "bench", NewLuaIncrementer(connect()))
}

func BenchmarkTxIncrementer(b *testing.B) {
	benchIncrementer(b, "bench", NewTxIncrementer(connect()))
}

func testIncrementer(t *testing.T, key string, incrementer Incrementer) {
	var expireSec int64 = 1
	var concurrent = 4
	var rounds = 100

	incr := func() int64 {
		v, err := incrementer.Incr(key, expireSec)
		if err != nil {
			panic(err)
		}
		return v
	}

	wg := sync.WaitGroup{}
	for c := 0; c < concurrent; c++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < rounds; i++ {
				incr()
			}
		}()
	}
	wg.Wait()

	if counter := incr(); counter != int64(concurrent*rounds)+1 {
		t.Errorf("counter should be %d, but got %d", concurrent*rounds+1, counter)
	}

	time.Sleep(time.Duration(expireSec) * time.Second)

	if counter := incr(); counter != 1 {
		t.Errorf("counter is not exipred after %d seconds", expireSec)
	}
}

func benchIncrementer(b *testing.B, key string, incrementer Incrementer) {
	for i := 0; i < b.N; i++ {
		_, err := incrementer.Incr(key, 60)
		if err != nil {
			panic(err)
		}
	}
}
