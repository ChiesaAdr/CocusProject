package cache

import (
	"cocus/internal/config"
	"log"
	"regexp"
	"sort"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

//Used to sort the cache
type pair struct {
	key   int
	value string
}

type pairList []pair

//Operations to sort
func (p pairList) Len() int           { return len(p) }
func (p pairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p pairList) Less(i, j int) bool { return p[i].key < p[j].key }

func (c *Cache) NewRedisClient() {

	//Use redisUrl and redisPassword to remote/cluster Redis.
	//To local is used the default(localhost:6379)
	redisUrl := config.GetRedisURL()
	redisPassword := config.GetRedisPassword()

	c.redis = redis.NewClient(&redis.Options{
		Addr:         redisUrl,
		Password:     redisPassword,
		DB:           0,
		OnConnect:    c.OnRedisConnect,
		DialTimeout:  1 * time.Minute,
		ReadTimeout:  1 * time.Minute,
		WriteTimeout: 1 * time.Minute,
	})
}

func (c *Cache) OnRedisConnect(cn *redis.Conn) error {
	return nil
}

//Monitoring a connection state with redis
//If you receive the LOG: Redis disconnected, please,
// check if the redis is UP and configured
func (c *Cache) RedisMonitor() {
	for {
		result, _ := c.redis.Ping().Result()
		if result == "PONG" {
			time.Sleep(60 * time.Second)
			continue
		}
		for {
			log.Println("core.RedisMonitor() - Redis disconnected")
			result, _ := c.redis.Ping().Result()
			if result == "PONG" {
				log.Println("core.RedisMonitor() - Redis reconnected")
				break
			}
			time.Sleep(10 * time.Second)
		}
	}
}

//Get value by key
func (c *Cache) RedisGet(key string) string {
	return c.redis.Get(key).Val()
}

//Delete the key in redis, don't look at the id to sort
func (c *Cache) RedisDel(key string) error {
	keys, _, err := c.redis.Scan(0, key+"&*", 0).Result()
	if err != nil {
		log.Printf("ERROR-Redis to del entry: %s", err)
		return nil
	}
	if len(keys) == 0 {
		log.Printf("ERROR-Redis to del entry: key nil")
		return nil
	}
	return c.redis.Del(keys[0]).Err()
}

//Add new message to redis
//TODO: Use ZADD to add in redis(Best way to sort)
func (c *Cache) RedisSet(key string, data string) error {
	//History is limited to 20 messages
	if c.RedisSize() > int64(c.limit) {
		c.RedisDelFristKey()
	}
	return c.redis.Set(key, data, 0).Err()
}

func (c *Cache) RedisGetAllKeys() []string {
	keys, _, err := c.redis.Scan(0, "*", 0).Result()
	if err != nil {
		log.Printf("ERROR-Redis to get all value: %s", err)
		return nil
	}
	return keys
}

//Get all the messages, but in the correct order
//used variable after `&` to sort and guarantee it
func (c *Cache) RedisGetAllValues() []string {
	var res []string = []string{""}
	p := make(pairList, 0)
	//Regex to find the order
	re1 := regexp.MustCompile(`&(\d+)`)
	iter := c.redis.Scan(0, "*", 0).Iterator()
	for iter.Next() {
		result := re1.FindStringSubmatch(iter.Val())
		if len(result) < 2 {
			continue
		}
		id, _ := strconv.Atoi(result[1])
		p = append(p, pair{key: id, value: iter.Val()})
	}
	//Sort struct by Key: id
	sort.Sort(p)
	for _, k := range p {
		res = append(res, c.RedisGet(k.value))
	}
	return res
}

func (c *Cache) RedisSize() int64 {
	return c.redis.DBSize().Val()
}

//Delete the oldest message
//TODO: Use LREM to del the oldest key in redis
func (c *Cache) RedisDelFristKey() error {
	p := make(pairList, 0)
	//Regex to find the order
	re1 := regexp.MustCompile(`&(\d+)`)
	iter := c.redis.Scan(0, "*", 0).Iterator()
	for iter.Next() {
		result := re1.FindStringSubmatch(iter.Val())
		if len(result) < 2 {
			continue
		}
		id, _ := strconv.Atoi(result[1])
		p = append(p, pair{key: id, value: iter.Val()})
	}
	//Sort struct by Key: id
	sort.Sort(p)
	for _, k := range p {
		//Delete de first entry
		return c.redis.Del(k.value).Err()
	}
	return nil
}

func (c *Cache) RedisFlushAll() error {
	return c.redis.FlushAll().Err()
}
