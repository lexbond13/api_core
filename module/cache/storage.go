package cache

import (
	"github.com/lexbond13/api_core/config"
	"time"

	"github.com/pkg/errors"
	"gopkg.in/redis.v5"
)

var client ICacheStorage

type Config struct {
	Network  string `json:"network"`
	Address  string `json:"address"`
	Password string `json:"password"`
	Database uint   `json:"database"`
	PoolSize int    `json:"pool_size" validate:"min=0"`
}

type ICacheStorage interface {
	Get(key string) (string, error)
	Set(key string, value interface{}, ttl int) error
	Del(key string) (bool, error)
}

func Init(configCache *config.Cache) error {
	// change cache service this!
	// check for enabled redis service
	if configCache.Redis.IsUsed {
		err := InitRedis(configCache.Redis)
		return err
	}
	return nil
}

func InitRedis(redisCfg *config.Redis) error {
	var err error
	client, err = NewRedisClient(redisCfg)
	if err != nil {
		return errors.Wrap(err, "Error init redis_env package")
	}

	return nil
}

func GetClient() ICacheStorage {
	return client
}

type RedisClient struct {
	*redis.Client
}

type RedisConfig struct {
	Network  string `json:"network"`
	Address  string `json:"address"`
	Password string `json:"password"`
	Database uint   `json:"database"`
	PoolSize int    `json:"pool_size" validate:"min=0"`
}

func NewRedisClient(redisCfg *config.Redis) (c *RedisClient, err error) {
	options := &redis.Options{
		Network:  redisCfg.Network,
		Addr:     redisCfg.Address,
		Password: redisCfg.Password,
		DB:       int(redisCfg.Database),
		PoolSize: int(redisCfg.PoolSize),
	}

	cli := redis.NewClient(options)

	_, err = cli.Ping().Result()
	if err != nil {
		err = errors.Wrap(err, "Can't connect to cache")
		return
	}
	c = &RedisClient{cli}
	return
}

func (r *RedisClient) Get(key string) (string, error) {
	result := r.Client.Get(key)

	// handle not found err
	if result.Err() == redis.Nil {
		return "", nil
	}

	return result.Result()
}

func (r *RedisClient) Set(key string, value interface{}, ttl int) error {
	duration := time.Duration(ttl) * time.Second
	result := r.Client.Set(key, value, duration)

	return result.Err()
}

func (r *RedisClient) Del(key string) (bool, error) {
	result := r.Client.Del(key)

	return result.Val() > 0, result.Err()
}
