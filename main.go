package main

import (
	"fmt"
	"github.com/1r0npipe/ethereum-parser/parser"
	"github.com/1r0npipe/ethereum-parser/storage"
	"github.com/sirupsen/logrus"
)

const (
	ethereumURL = "https://cloudflare-eth.com"
	id          = 1
	address     = "0x0000000000000000000000000000000000000011" // Example address
	blockRange  = 1
)

func main() {
	logrus.Info("Starting Ethereum Parser App")

	// Initialize in-memory storage
	memStorage := storage.NewMemoryStorage()

	// Create a new parser with the Ethereum JSON-RPC URL and memory storage
	ethereumParser := parser.NewEthereumParser(ethereumURL, memStorage)

	logrus.Info("App Initialized")

	// Example how we do subscription to the address
	if memStorage.Subscribe(address) {
		logrus.Infof("Subscribed to address: %s", address)
	} else {
		logrus.Warnf("Address %s is already subscribed", address)
	}

	// Getting transactions for the address
	transactions, err := ethereumParser.GetTransactions(address, blockRange, id)
	if err != nil {
		logrus.Errorf("Failed to get transactions: %v", err)
		return
	}

	if len(transactions) == 0 {
		logrus.Infof("No transactions found for address %s", address)
	} else {
		for _, tx := range transactions {
			logrus.Infof("Transaction: Hash=%s From=%s To=%s Block=%s", tx.Hash, tx.From, tx.To, tx.BlockNumber)
		}
	}

	// Example retrieval of all subscribed addresses for example purposes
	subscribedAddresses := memStorage.GetAddresses()
	fmt.Println("Subscribed Addresses:")
	for _, addr := range subscribedAddresses {
		fmt.Println(addr)
	}
}
