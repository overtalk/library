## Golang 操作 Redis 的基本方法

### 摘要
- 看到有位老哥写了PHP操作Redis的基本方法，于是就有了这篇博客
- 话不多说，直接上代码

### Redis 链接池
- Redis有集群版和单机版，当配置中的`Addrs`写成`ip:port`的数组形式即可链接集群版Redis
```Go
package redis

import (
	"errors"

	"github.com/go-redis/redis"
)

type Cfg struct {
	Addrs    []string `json:"redis_addrs" env:"REDIS_ADDRS" envDefault:"127.0.0.1:6379"`
	Pwd      string   `json:"redis_pwd" env:"REDIS_PWD"`
	PoolSize int      `json:"redis_pool_size" env:"REDIS_POOL_SIZE" envDefault:"1000"`
	DB       int      `json:"redis_db" env:"REDIS_DB"` // 单机模式下选择使用哪个DB，集群模式下无效
}

func (c Cfg) Connect() (redis.Cmdable, error) {
	addrNum := len(c.Addrs)
	if addrNum == 0 {
		return nil, errors.New("redis addr is absent")
	}

	if addrNum > 1 {
		return redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    c.Addrs,
			Password: c.Pwd,
			PoolSize: c.PoolSize,
		}), nil
	}
	return redis.NewClient(&redis.Options{
		Addr:     c.Addrs[0],
		Password: c.Pwd,
		PoolSize: c.PoolSize,
		DB:       c.DB,
	}), nil
}
``` 

- 接下来是一份测试，所有配置通过环境变量的形式来读取
```go
package redis

import (
	"os"
	"testing"
	"time"

	"github.com/caarlos0/env"
)

func TestConnect(t *testing.T) {
	os.Setenv("REDIS_ADDRS", "127.0.0.1:6379")
	os.Setenv("REDIS_PWD", "")
	os.Setenv("REDIS_POOL_SIZE", "100")
	os.Setenv("REDIS_DB", "1")

	c := Cfg{}
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
```
### Redis 的数据操作
- 说完了如何链接到 Redis，接下来就说一说对 Redis 进行数据操作吧
- Redis 中的数据结构可以参考官网（https://redis.io/commands）
- 这儿只举例常用的一些数据结构，其他的读者可以自行研究，上文中使用的Redis库（github.com/go-redis/redis"）的 redis.Cmdable 已经实现了所有的数据结构，有兴趣可以自己去看
- 最常用的 Get & Set 在上一个测试中已经有了，就不多做叙述

### hash
```go
package redis

import (
	"os"
	"testing"
	"time"

	"github.com/caarlos0/env"
)

func TestRedis_Hash(t *testing.T) {
	os.Setenv("REDIS_ADDRS", "127.0.0.1:6379")
	os.Setenv("REDIS_PWD", "")
	os.Setenv("REDIS_POOL_SIZE", "100")
	os.Setenv("REDIS_DB", "1")

	c := Cfg{}
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
```
- 测试运行结果如下


### sorted sets
- 排序的集合
- 集合里面是一对对的 k-v，并且v为float64，根据v进行排序
- 这个数据结构可以用来做排行榜

```go
package redis

import (
	"os"
	"testing"
	"time"

	"github.com/caarlos0/env"
	"github.com/go-redis/redis"
)

func TestRedis_Sorted_Sets(t *testing.T) {
	os.Setenv("REDIS_ADDRS", "127.0.0.1:6379")
	os.Setenv("REDIS_PWD", "")
	os.Setenv("REDIS_POOL_SIZE", "100")
	os.Setenv("REDIS_DB", "1")

	c := Cfg{}
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
```
- 测试运行结果如下

### Lists
- 列表

```go
package redis

import (
	"os"
	"testing"
	"time"

	"github.com/caarlos0/env"
)

func TestRedis_Lists(t *testing.T) {
	os.Setenv("REDIS_ADDRS", "127.0.0.1:6379")
	os.Setenv("REDIS_PWD", "")
	os.Setenv("REDIS_POOL_SIZE", "100")
	os.Setenv("REDIS_DB", "1")

	c := Cfg{}
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
```
- 测试运行结果如下
