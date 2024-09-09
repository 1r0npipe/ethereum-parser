package storage

import "sync"

// Storage defines the interface for a storage system.
type Storage interface {
	// Subscribe adds an address to the storage and returns true if added.
	Subscribe(address string) bool

	// GetAddresses returns all subscribed addresses.
	GetAddresses() []string
}

// MemoryStorage implements the Storage interface using an in-memory map.
type MemoryStorage struct {
	addresses map[string]interface{}
	mu        sync.Mutex
}

// NewMemoryStorage creates and returns a new MemoryStorage instance.
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		addresses: make(map[string]interface{}),
	}
}

// Subscribe adds an address to the memory storage and returns true if it was added.
func (ms *MemoryStorage) Subscribe(address string) bool {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	if _, exists := ms.addresses[address]; exists {
		return false // Address is already subscribed
	}
	// just initialize but can keep anything
	ms.addresses[address] = struct{}{}
	return true
}

// GetAddresses returns all subscribed addresses.
func (ms *MemoryStorage) GetAddresses() []string {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	var addresses []string
	for address := range ms.addresses {
		addresses = append(addresses, address)
	}
	return addresses
}
