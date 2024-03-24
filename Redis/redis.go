package redispkg

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

type RedisStruct struct {
	client *redis.Client
}

func NewRedisClient() *RedisStruct {
	redisURL := os.Getenv("redisUrl")
	redisPass := os.Getenv("redisPwd")
	client := redis.NewClient(&redis.Options{
		Addr:     redisURL,
		Password: redisPass, // no password set
		DB:       0,         // use default DB
	})

	pong, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Panic("Could not connect to Redis\n", err)
		return nil
	}
	log.Println(pong) // Output: PONG
	return &RedisStruct{
		client: client,
	}
}

func (R *RedisStruct) GetObj(username string) []byte {
	// Get User From Redis
	ctx := context.Background()
	val, err := R.client.Get(ctx, username).Bytes()
	if err != nil {
		log.Print("user not found in redis. ", err)
		return nil
	}

	// R.client.`Set(ctx, username, username+"1", 0)
	// var userV`al any
	// Convert the []Byte to a struct
	return val
	// json.NewDecoder(val).Decode()
	// userVal := helpers.ConvertFromJSON([]byte(val))
	// if userVal == nil {
	// 	return nil
	// }

	// return userVal

}

func (R *RedisStruct) GetUsers() []string {
	ctx := context.Background()
	keys := R.client.Keys(ctx, "*").Val()
	var jsonArray []string
	for _, key := range keys {
		jsonStr := R.client.Get(ctx, key).Val()
		jsonArray = append(jsonArray, jsonStr)
	}
	return jsonArray
}

func (R *RedisStruct) SetObj(key string, user any) {
	ctx := context.Background()
	userJSON, err := json.Marshal(user)
	if err != nil {
		log.Println("Error marshaling user struct to JSON:", err)
		return
	}
	R.client.Set(ctx, key, userJSON, 0)
	// Convert the user to a byteArray
	// userByteArray := helpers.ConvertToJSON(user)
	// if userByteArray == nil {
	// 	return
	// }
	// log.Println("json to be stored:", userByteArray)
	// var a any
	// json.Unmarshal(userByteArray, &a)
	// log.Println(a)
	// R.client.Set(ctx, key, userByteArray, 0)

}
