package store

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"habr-searcher/internal/model"
)

type Store struct {
	rdb *redis.Client
	ctx context.Context
}

func New(ctx context.Context) *Store {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	return &Store{
		rdb: rdb,
		ctx: ctx,
	}
}
func (s *Store) Set(key string, value interface{}) error {
	v, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return s.rdb.Set(s.ctx, key, v, 0).Err()
}

func (s *Store) GetUser(key string) (model.User, error) {
	b, err := s.rdb.Get(s.ctx, key).Bytes()
	u := model.User{}
	if err != nil {
		return u, err
	}
	err = json.Unmarshal(b, &u)
	return u, err
}
