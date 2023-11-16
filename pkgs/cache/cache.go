package cache

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/nitishm/go-rejson/v4"
	"github.com/redis/go-redis/v9"
)

var ErrKeyNotFound = errors.New("not found")

type RedsiStore struct {
	*redis.Client
	*rejson.Handler
}

type Object interface {
	Key() string
}

type JSONObject interface {
	Object
	New() JSONObject
	// json.Marshaler
	// json.Unmarshaler
}

func NewRedsiStore(ctx context.Context, opts *redis.Options) *RedsiStore {
	client := redis.NewClient(opts)
	handler := rejson.NewReJSONHandler()
	handler.SetGoRedisClientWithContext(ctx, client)

	return &RedsiStore{
		Client:  client,
		Handler: handler,
	}
}

func (rs *RedsiStore) ObjectGet(ctx context.Context, object Object) *redis.StringCmd {
	return rs.Client.Get(ctx, object.Key())
}

func (rs *RedsiStore) ObjectSet(ctx context.Context, object Object, expiration time.Duration) error {
	return rs.Client.Set(ctx, object.Key(), object, expiration).Err()
}

func (rs *RedsiStore) ObjectDel(ctx context.Context, object Object) error {
	return rs.Client.Del(ctx, object.Key()).Err()
}

func (rs *RedsiStore) JSONObjectGet(jsonObject JSONObject, path string) error {
	res, err := rs.Handler.JSONGet(jsonObject.Key(), path)
	if err != nil {
		return err
	}

	if res == redis.Nil {
		return ErrKeyNotFound
	}

	return json.Unmarshal(res.([]byte), jsonObject)
}

func (rs *RedsiStore) JSONObjectSet(jsonObject JSONObject, path string, expiration time.Duration) error {
	_, err := rs.Handler.JSONSet(jsonObject.Key(), path, jsonObject)
	if err != nil {
		return err
	}
	return rs.Client.Expire(context.Background(), jsonObject.Key(), expiration).Err()
}

func (rs *RedsiStore) JSONObjectDel(jsonObject JSONObject, path string) error {
	_, err := rs.Handler.JSONDel(jsonObject.Key(), path)
	return err
}
