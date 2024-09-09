package parser

import (
	"encoding/json"
	"fmt"
	"github.com/1r0npipe/ethereum-parser/storage"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strings"
)

// Transaction represents a transaction in the Ethereum blockchain.
type Transaction struct {
	Hash        string `json:"hash"`
	From        string `json:"from"`
	To          string `json:"to"`
	BlockNumber string `json:"blockNumber"`
}

// Parser is an interface defining the Ethereum blockchain parser methods.
type Parser interface {
	GetCurrentBlock(id int) (int, error)
	Subscribe(address string) bool
	GetTransactions(address string, blockRange int, id int) ([]Transaction, error)
}

// EthereumParser implements the Parser interface.
type ethereumParser struct {
	URL     string
	storage storage.Storage
}

// NewEthereumParser initializes a new Ethereum Parser instance.
func NewEthereumParser(URL string, storage storage.Storage) Parser {
	return &ethereumParser{
		URL:     URL,
		storage: storage,
	}
}

// GetCurrentBlock fetches the latest block number from the Ethereum network.
func (p *ethereumParser) GetCurrentBlock(id int) (int, error) {
	requestBody := fmt.Sprintf(`{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":%d}`, id)

	resp, err := http.Post(p.URL, "application/json", strings.NewReader(requestBody))
	if err != nil {
		return -1, fmt.Errorf("failed to fetch block number: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return -1, fmt.Errorf("failed to read response body: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return -1, fmt.Errorf("failed to parse JSON response: %v", err)
	}

	blockNumberHex, ok := result["result"].(string)
	if !ok {
		return -1, fmt.Errorf("unexpected result format: %v", result)
	}

	var blockNumber int64
	_, err = fmt.Sscanf(blockNumberHex, "0x%x", &blockNumber)
	if err != nil {
		return -1, fmt.Errorf("error parsing block number: %v", err)
	}

	return int(blockNumber), nil
}

// Subscribe adds an address to the memory storage.
func (p *ethereumParser) Subscribe(address string) bool {
	return p.storage.Subscribe(address)
}

// GetTransactions gets all transactions for a given address by iterating over blocks range (not sure if I understood correct here?).
func (p *ethereumParser) GetTransactions(address string, blockRange int, id int) ([]Transaction, error) {
	var transactions []Transaction
	latestBlock, err := p.GetCurrentBlock(id)
	if err != nil {
		return nil, err
	}

	for blockNum := latestBlock; blockNum > latestBlock-blockRange; blockNum-- {
		blockTransactions, err := p.getBlockTransactions(blockNum, id)
		if err != nil {
			logrus.Errorf("Failed to fetch transactions for block %d: %v", blockNum, err)
			continue
		}

		for _, tx := range blockTransactions {
			if tx.From == address || tx.To == address {
				transactions = append(transactions, tx)
			}
		}
	}

	return transactions, nil
}

// getBlockTransactions gets all transactions from a block.
func (p *ethereumParser) getBlockTransactions(blockNumber int, id int) ([]Transaction, error) {
	requestBody := fmt.Sprintf(`{"jsonrpc":"2.0","method":"eth_getBlockByNumber","params":["0x%x", true],"id":%d}`, blockNumber, id)

	resp, err := http.Post(p.URL, "application/json", strings.NewReader(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch block %d: %v", blockNumber, err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to parse block JSON response: %v", err)
	}

	block, ok := result["result"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid block data structure")
	}

	transactionsData, ok := block["transactions"].([]interface{})
	if !ok {
		logrus.Warnf("No transactions found for block %d", blockNumber)
		return nil, nil
	}

	var transactions []Transaction
	for _, txData := range transactionsData {
		txMap, ok := txData.(map[string]interface{})
		if !ok {
			continue
		}

		tx := Transaction{
			Hash:        txMap["hash"].(string),
			From:        txMap["from"].(string),
			To:          txMap["to"].(string),
			BlockNumber: fmt.Sprintf("0x%x", blockNumber),
		}

		transactions = append(transactions, tx)
	}

	return transactions, nil
}
