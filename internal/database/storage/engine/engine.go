package engine

import "sync"

type Engine struct {
	mutex sync.RWMutex
	data  map[string]string
}

func NewEngine() *Engine {
	return &Engine{
		data: make(map[string]string),
	}
}

func (e *Engine) Set(key, value string) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	e.data[key] = value
}

func (e *Engine) Get(key string) (string, bool) {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	value, ok := e.data[key]

	return value, ok
}

func (e *Engine) Del(key string) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	delete(e.data, key)
}
