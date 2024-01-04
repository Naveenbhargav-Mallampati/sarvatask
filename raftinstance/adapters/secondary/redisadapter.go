package secondary

import (
	"fmt"

	"github.com/go-redis/redis"
)

func GetData(key string) string {
	Rediscon := *redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	value, err := Rediscon.Get(key).Result()
	fmt.Println(key)
	fmt.Println(value)
	if err != nil {
		fmt.Println("Error getting data")
	}

	return value

}

func SetData(key, value string) bool {
	Rediscon := *redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	err := Rediscon.Set(key, value, 0).Err()
	if err != nil {
		return false
	}
	return true
}
