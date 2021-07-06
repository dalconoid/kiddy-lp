package storage

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"strconv"
)

const (
	baseballKey = "baseball"
	footballKey = "football"
	soccerKey = "soccer"
)

// RedisStorage is a redis storage
type RedisStorage struct {
	Client *redis.Client
}

// NewRedisStorage creates redis storage on addr with pass
func NewRedisStorage(addr, pass string) (*RedisStorage, error) {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
		Password: pass,
		DB: 0,
	})

	err := client.Ping(context.Background()).Err()
	if err != nil {
		return nil, err
	}

	return &RedisStorage{Client: client}, nil
}

// WriteLineRate writes rate k of a line "line" to storage
func (s *RedisStorage) WriteLineRate(k float64, line string) error {
	if err := s.Client.Set(context.Background(), line, k, 0).Err(); err != nil {
		return err
	}
	return nil
}

// GetLineRate gets rate of a line "line" from storage
func (s *RedisStorage) GetLineRate(line string) (float64, error) {
	v, err := s.Client.Get(context.Background(), line).Result()
	if err != nil {
		return 0, err
	}
	k, _ := strconv.ParseFloat(v, 64)
	return k, nil
}

// CheckConnection checks connection to storage and checks if line rates are already present in storage
func (s *RedisStorage) CheckConnection() error {
	if err := s.Client.Ping(context.Background()).Err(); err != nil {
		return err
	}

	keys := make([]string, 0, 3)

	if err := s.Client.Get(context.Background(), baseballKey).Err(); err == redis.Nil {
		keys = append(keys, baseballKey)
	}
	if err := s.Client.Get(context.Background(), footballKey).Err(); err == redis.Nil {
		keys = append(keys, footballKey)
	}
	if err := s.Client.Get(context.Background(), soccerKey).Err(); err == redis.Nil {
		keys = append(keys, soccerKey)
	}

	if len(keys) > 0 {
		return fmt.Errorf("connection to lines %s hs not been established yet", keys)
	}

	return nil
}