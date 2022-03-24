package tests

import (
	"cocus/internal/cache"
	"strings"
	"testing"
)

func TestBasicRedis(t *testing.T) {
	key := "1234"
	data := "asdf"
	c := cache.NewCache()

	err := c.GetRedis().Set(key, data, 0).Err()

	if err != nil {
		t.Fatalf("ERROR to set Redis")
	}

	value := c.GetRedis().Get(key)
	if value.Err() != nil {
		t.Fatalf("ERROR to Get Redis")
	}

	if strings.Compare(value.Val(), data) != 0 {
		t.Error("ERROR Value different")
	}
}

func TestFlushRedis(t *testing.T) {

	c := cache.NewCache()
	err := c.GetRedis().FlushAll().Err()
	if err != nil {
		t.Error("ERROR to Flush Redis")
	}
}

func TestGetAllRedisCheckSize(t *testing.T) {
	key := "1234"
	data := "asdf"
	c := cache.NewCache()

	err := c.GetRedis().Set(key, data, 0).Err()
	if err != nil {
		t.Fatalf("ERROR to set Redis")
	}

	var res []string
	iter := c.GetRedis().Scan(0, "*", 0).Iterator()
	for iter.Next() {
		res = append(res, c.RedisGet(iter.Val()))
	}

	n := c.GetRedis().DBSize().Val()
	if len(res) != int(n) {
		t.Error("ERROR to read all the values!")
	}
}

func TestDelFristEntry(t *testing.T) {

	c := cache.NewCache()

	n1 := c.GetRedis().DBSize().Val()

	iter := c.GetRedis().Scan(0, "*", 0).Iterator()
	if iter.Next() {
		c.GetRedis().Del(iter.Val())
	}

	n2 := c.GetRedis().DBSize().Val()
	if n1 == n2 {
		t.Error("ERROR to read all the values!")
	}
}
