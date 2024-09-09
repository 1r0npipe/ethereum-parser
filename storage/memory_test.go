package storage

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMemoryStorage_Subscribe(t *testing.T) {
	ms := NewMemoryStorage()

	// Test adding a new address
	address := "0x0000000000000000000000000000000000000011"
	added := ms.Subscribe(address)
	assert.True(t, added, "Expected address to be added")

	// Test adding the same address again
	added = ms.Subscribe(address)
	assert.False(t, added, "Expected address not to be added again")
}

func TestMemoryStorage_GetAddresses(t *testing.T) {
	ms := NewMemoryStorage()

	// Test subscribing to addresses
	addresses := []string{
		"0x0000000000000000000000000000000000000011",
		"0x0000000000000000000000000000000000000012",
	}

	for _, address := range addresses {
		ms.Subscribe(address)
	}

	// Test retrieval of subscribed addresses
	retrievedAddresses := ms.GetAddresses()
	assert.ElementsMatch(t, addresses, retrievedAddresses, "Retrieved addresses do not match")
}
