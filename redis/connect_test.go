package redis_test

import (
	"os"
	"testing"
	"time"

	"github.com/caarlos0/env"
	"github.com/go-redis/redis"

	. "web-layout/utils/redis"
)

func TestConnect(t *testing.T) {
	os.Setenv("REDIS_ADDRS", "127.0.0.1:6379")
	os.Setenv("REDIS_PWD", "")
	os.Setenv("REDIS_POOL_SIZE", "100")
	os.Setenv("REDIS_DB", "1")

	c := Config{}
	if err := env.Parse(&c); err != nil {
		t.Error(err)
		return
	}

	redisPool, err := c.Connect()
	if err != nil {
		t.Error(err)
		return
	}

	res, err := redisPool.Set("test_key1", "test_value", time.Minute*10).Result()
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(res)
}

func TestRedis_Hash(t *testing.T) {
	os.Setenv("REDIS_ADDRS", "127.0.0.1:6379")
	os.Setenv("REDIS_PWD", "")
	os.Setenv("REDIS_POOL_SIZE", "100")
	os.Setenv("REDIS_DB", "1")

	c := Config{}
	if err := env.Parse(&c); err != nil {
		t.Error(err)
		return
	}

	redisPool, err := c.Connect()
	if err != nil {
		t.Error(err)
		return
	}

	hashKey := "HashKey"
	data := map[string]interface{}{
		"field1": "value1",
		"field2": 2,
		"field3": true,
	}
	for k, v := range data {
		_, err := redisPool.HSet(hashKey, k, v).Result()
		if err != nil {
			t.Error(err)
			return
		}
	}

	// get all
	results, err := redisPool.HGetAll(hashKey).Result()
	if err != nil {
		t.Error(err)
		return
	}

	t.Log("HGetAll Results : ", results)

	// get certain field
	value, err := redisPool.HGet(hashKey, "field1").Result()
	if err != nil {
		t.Error(err)
		return
	}

	t.Log("HGet field1 : ", value)
}

func TestRedis_Sorted_Sets(t *testing.T) {
	os.Setenv("REDIS_ADDRS", "127.0.0.1:6379")
	os.Setenv("REDIS_PWD", "")
	os.Setenv("REDIS_POOL_SIZE", "100")
	os.Setenv("REDIS_DB", "1")

	c := Config{}
	if err := env.Parse(&c); err != nil {
		t.Error(err)
		return
	}

	redisPool, err := c.Connect()
	if err != nil {
		t.Error(err)
		return
	}

	sortedSetKey := "SortedSets"
	data := map[interface{}]float64{
		"player1": 12,
		"player2": 2,
		"player3": 3,
	}

	// 设置两遍
	for i := 0; i < 2; i++ {
		for k, v := range data {
			_, err := redisPool.ZIncr(sortedSetKey, redis.Z{
				Member: k,
				Score:  v,
			}).Result()
			if err != nil {
				t.Error(err)
				return
			}
		}
	}

	// 根据排名获取数据
	// 有多种获取数据的方法
	results, err := redisPool.ZRevRangeWithScores(sortedSetKey, 0, 10).Result()
	if err != nil {
		t.Error(err)
		return
	}

	for _, v := range results {
		t.Log(v)
	}
}

func TestRedis_Lists(t *testing.T) {
	os.Setenv("REDIS_ADDRS", "127.0.0.1:6379")
	os.Setenv("REDIS_PWD", "")
	os.Setenv("REDIS_POOL_SIZE", "100")
	os.Setenv("REDIS_DB", "1")

	c := Config{}
	if err := env.Parse(&c); err != nil {
		t.Error(err)
		return
	}

	redisPool, err := c.Connect()
	if err != nil {
		t.Error(err)
		return
	}

	// 向list中插入数据
	listKey := "List"
	data := []interface{}{
		true,
		12,
		"test",
		22,
		"player3",
		3.9,
	}

	_, err = redisPool.LPush(listKey, data...).Result()
	if err != nil {
		t.Error(err)
		return
	}

	// 获取长度
	len, err := redisPool.LLen(listKey).Result()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("LLen = ", len)

	// pop
	popData, err := redisPool.LPop(listKey).Result()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("LPop = ", popData)
}
