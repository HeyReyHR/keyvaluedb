package storage

import "errors"

var (
	ErrNotFound = errors.New("not found")
)

type Engine interface {
	Set(key, value string)
	Get(key string) (string, bool)
	Del(key string)
}
type Storage struct {
	engine Engine
}

func NewStorage(engine Engine) *Storage {
	return &Storage{
		engine: engine,
	}
}

func (s *Storage) Set(key, value string) error {
	s.engine.Set(key, value)
	return nil
}

func (s *Storage) Get(key string) (string, error) {
	val, ok := s.engine.Get(key)
	if !ok {
		return "", ErrNotFound
	}
	return val, nil
}

func (s *Storage) Del(key string) error {
	s.engine.Del(key)
	return nil
}
