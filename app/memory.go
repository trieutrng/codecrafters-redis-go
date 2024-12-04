package main

type Memory struct {
	store  map[string]string
	expiry <-chan string
}

func NewMemory() *Memory {
	memory := &Memory{
		store:  make(map[string]string),
		expiry: make(<-chan string),
	}

	// watch expiry event asynchronously
	go memory.expiryWatcher()

	return memory
}

func (m *Memory) Get(key string) string {
	val, ok := m.store[key]
	if !ok {
		return ""
	}
	return val
}

func (m *Memory) Put(key, val string) {
	m.store[key] = val
}

func (m *Memory) expiryWatcher() {
	for expiredKey := range m.expiry {
		delete(m.store, expiredKey)
	}
}
