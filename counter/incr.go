package counter

import (
	"time"

	"github.com/go-redis/redis"
)

var incrScript *redis.Script

func init() {
	incrScript = redis.NewScript(`local c
c = redis.call("incr",KEYS[1])
if tonumber(c) == 1 then
    redis.call("expire",KEYS[1],ARGV[1])
end
return c`)
}

type Incrementer interface {
	Incr(key string, expireSec int64) (count int64, err error)
}

func NewLuaIncrementer(pool redis.UniversalClient) Incrementer {
	return &LuaIncrementer{pool}
}

func NewTxIncrementer(pool redis.UniversalClient) Incrementer {
	return &TxIncrementer{pool}
}

type LuaIncrementer struct {
	pool redis.UniversalClient
}

func (i *LuaIncrementer) Incr(key string, expireSec int64) (count int64, err error) {
	v, err := incrScript.Run(i.pool, []string{key}, expireSec).Result()
	return v.(int64), err
}

type TxIncrementer struct {
	pool redis.UniversalClient
}

func (i *TxIncrementer) Incr(key string, expireSec int64) (count int64, err error) {
	pipe := i.pool.TxPipeline()
	pipe.SetNX(key, 0, time.Duration(expireSec)*time.Second)
	incr := pipe.Incr(key)

	_, err = pipe.Exec()

	return incr.Val(), err
}
