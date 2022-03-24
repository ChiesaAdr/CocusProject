package tests

// func TestBasicRedis(t *testing.T) {
// 	key := "1234"
// 	data := "asdf"
// 	c := cache.NewCache()

// 	err := c.redis.Set(key, data, 0).Err()

// 	if err != nil {
// 		t.Fatalf("ERROR to set Redis")
// 	}

// 	value := c.redis.Get(key)
// 	if value.Err() != nil {
// 		t.Fatalf("ERROR to Get Redis")
// 	}

// 	if strings.Compare(value.Val(), data) != 0 {
// 		t.Error("ERROR Value different")
// 	}
// }

// func TestBasicRedis2(t *testing.T) {
// 	key := "12342367"
// 	data := "asdf"
// 	c := cache.NewCache()

// 	c.redis.Set(key, data, 0).Err()

// }

// func TestFlushRedis(t *testing.T) {

// 	c := cache.NewCache()
// 	err := c.redis.FlushAll().Err()
// 	if err != nil {
// 		t.Error("ERROR to Flush Redis")
// 	}
// }

// func TestGetAllRedis(t *testing.T) {

// 	c := cache.NewCache()

// 	// iter := c.redis.Scan(0, "*", 0).Iterator()
// 	// for iter.Next() {
// 	// 	t.Errorf("keys %v", iter.Val())
// 	// 	// c.redis.Del(iter.Val())
// 	// }
// 	// var res []string
// 	iter := c.redis.Scan(0, "*", 0).Iterator()
// 	for iter.Next() {
// 		t.Errorf("val %s", iter.Val())
// 		// res = append(res, c.RedisGet(iter.Val()))
// 	}
// 	// keys := c.RedisGetAllValues()
// 	// keys, _, err := c.redis.Scan(0, "*", 0).Result()
// 	// t.Errorf("keys %s", res)
// 	// if err != nil {
// 	// t.Errorf("Val %s  err %v", err.String(), err)
// 	// }

// 	n := c.redis.DBSize().Val()
// 	t.Errorf("keys %v ", n)
// }

// func TestDelFristEntry(t *testing.T) {

// 	c := cache.NewCache()

// 	keys, _, err := c.redis.Scan(0, "*", 0).Result()
// 	t.Errorf("keys %v err %v", keys, err)
// 	// if err != nil {
// 	// t.Errorf("Val %s  err %v", err.String(), err)
// 	// }

// 	n := c.redis.DBSize().Val()
// 	t.Errorf("keys %v err %v", n, err)

// 	// c.redis.LPop("*")
// 	iter := c.redis.Scan(0, "*", 0).Iterator()
// 	if iter.Next() {
// 		t.Errorf("keys %v", iter.Val())
// 		c.redis.Del(iter.Val())
// 	}

// 	keys, _, err = c.redis.Scan(0, "*", 0).Result()
// 	t.Errorf("keys %v err %v", keys, err)

// 	n = c.redis.DBSize().Val()
// 	t.Errorf("keys %v err %v", n, err)
// }

// func getWinnerPlayer() string {
// 	players := []string{"Mohammad", "Ali", "John", "Abdullah", "Farida"}
// 	return players[rand.Intn(len(players))]
// }

// func TestAddSorted(t *testing.T) {
// 	// key := "123456"
// 	// data := "asdf"
// 	c := cache.NewCache()

// 	const key = "players"

// 	players := []string{"Mohammad", "Ali", "John", "Abdullah", "Farida"}

// 	for i := 0; i < 5; i++ {
// 		// player := getWinnerPlayer()

// 		err := c.redis.ZIncr(players[i], redis.Z{
// 			Score:  1,
// 			Member: players[i],
// 		}).Err()

// 		if err != nil {
// 			t.Error(err)
// 		}
// 	}
// 	for i := 0; i < 5; i++ {
// 		result := c.redis.ZRevRangeWithScores(players[i], 0, -1).Val()

// 		t.Error(result)
// 	}
// 	// i, err := c.redis.ZAdd("asdf", redis.Z{2, key}).Result()
// 	// t.Errorf("i %v err %v", i, err)

// }
