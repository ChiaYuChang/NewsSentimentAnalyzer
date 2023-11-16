package cache_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/cache"
	"github.com/joho/godotenv"
	"github.com/nitishm/go-rejson/v4"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

func TestNotFound(t *testing.T) {
	err := godotenv.Load("../../.env")
	if err != nil {
		t.Log("skip test, .env file not found")
		return
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		t.Log("skip test, failed to connect to redis")
		return
	}

	ctx := context.Background()
	cmd := rdb.Get(ctx, "not-found")
	require.Equal(t, redis.Nil, cmd.Err())
}

type Employee struct {
	ID   int  `json:"id"`
	Name Name `json:"name"`
	Age  int  `json:"age"`
}

func NewEmployee(id int) *Employee {
	return &Employee{ID: id}
}

func (e *Employee) WithName(first, middle, last string) *Employee {
	e.Name.First = first
	e.Name.Middle = middle
	e.Name.Last = last
	return e
}

func (e *Employee) WithAge(age int) *Employee {
	e.Age = age
	return e
}

type Name struct {
	First  string `json:"first"`
	Last   string `json:"last,omitempty"`
	Middle string `json:"middle"`
}

func (e Employee) Key() string {
	return fmt.Sprintf("Employee-No.%05d", e.ID)
}

func (e Employee) New() cache.JSONObject {
	return &Employee{}
}

func TestJSONObject(t *testing.T) {
	err := godotenv.Load("../../.env")
	if err != nil {
		t.Log("skip test, .env file not found")
		return
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		t.Log("skip test, failed to connect to redis")
		return
	}

	h := rejson.NewReJSONHandler()
	h.SetGoRedisClientWithContext(context.Background(), rdb)
	store := cache.RedsiStore{
		Client:  rdb,
		Handler: h,
	}

	employee0 := NewEmployee(1).WithAge(30).WithName("John", "", "Doe")
	err = store.JSONObjectSet(employee0, ".", 10*time.Minute)
	require.NoError(t, err)

	t.Log(employee0.Key())

	employee1 := NewEmployee(1)
	err = store.JSONObjectGet(employee1, ".")
	require.NoError(t, err)

	t.Log(employee1)

}
