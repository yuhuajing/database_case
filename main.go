package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"os"
	"time"
)

var client *redis.Client

func init() {
	loadenv()
	connredis()
}
func connredis() {
	rpc := Env("REDIS_HOSTS", "")
	client = redis.NewClient(&redis.Options{Addr: rpc})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		fmt.Printf("Redis链接失败：%s\n", err)
		return
	}
	//fmt.Println(pong)
}
func loadenv() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	} // The Original .env
	_ = godotenv.Overload(".env")
}

func Env(key string, d string) string {
	env, ok := os.LookupEnv(key)
	if ok {
		return env
	}

	return d
}

func main() {
	key := "Alice"
	expiration := 15 * time.Second
	//value := "Elison"
	value, err := json.Marshal(Auther{Name: key, Age: 29})
	if err != nil {
		fmt.Printf("Json Marshal 编码数据错误：%s\n", err)
		return
	}
	err = setvalue(client, key, string(value), expiration)
	if err != nil {
		fmt.Printf("Redis Set键值数据错误：%s\n", err)
		return
	}
	values, err := getvalue(client, key)
	if err != nil {
		fmt.Printf("Redis Get键值数据错误：%s\n", err)
		return
	}
	valuestr := fmt.Sprintf("%s", values)
	var author Auther
	json.Unmarshal([]byte(valuestr), &author)
	fmt.Println(values)
}

func setvalue(c *redis.Client, key string, value interface{}, expiration time.Duration) error {
	err := c.Set(context.Background(), key, value, expiration).Err()
	return err
}

func getvalue(c *redis.Client, key string) (interface{}, error) {
	val, err := c.Get(context.Background(), key).Result()
	return val, err
}

type Auther struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}
