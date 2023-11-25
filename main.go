package main

import (
	"fmt"
	"log"
)

type Cacher interface {
	Get(int) (string, bool)
	Remove(int) error
	Set(int, string) error
}

type NopCache struct{}

func (c *NopCache) Get(int) (string, bool) {
	return "", false
}

func (c *NopCache) Remove(int) error {
	return nil
}

func (c *NopCache) Set(int, string) error {
	return nil
}

type Store struct {
	data  map[int]string
	cache Cacher
}

func NewStore(c Cacher) *Store {
	data := map[int]string{
		1: "Elon Musk is the new owner of Twitter",
		2: "Foo is not bar and bar is not baz",
		3: "Must watch AnthonyGG",
	}
	return &Store{
		data:  data,
		cache: c,
	}
}

func (s *Store) Get(key int) (string, error) {
	val, ok := s.cache.Get(key)
	if ok {
		if err := s.cache.Remove(key); err != nil {
			fmt.Println(err)
		}
		return val, nil
	}
	val, ok = s.data[key]
	if !ok {
		return "", fmt.Errorf("key not found: %d", key)
	}
	fmt.Println("returning key from internal storage")

	return val, nil
}

func (s *Store) getFromCache(key int) (string, bool) {
	val, ok := s.cache.Get(key)
	if ok {
		fmt.Println("returning key from cache")
		return val, ok
	}
	return "", false
}

func main() {
	// client := redis.NewClient(
	// 	&redis.Options{
	// 		Addr:     "localhost:6379",
	// 		Password: "",
	// 		DB:       0,
	// 	},
	// )

	s := NewStore(&NopCache{})

	for q := 0; q < 10; q++ {
		val, err := s.Get(1)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(val)
	}

}
