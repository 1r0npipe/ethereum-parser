# Ethereum Parser App

## Overview
This is a simple, and rather example of the Ethereum blockchain parser implemented in Go. It uses Ethereum's JSON-RPC interface to interact with the blockchain and subscribes to addresses for incoming and outgoing transactions.

## Requirements
- Go 1.20+
- logrus for logging
- mockery for mocking interfaces in tests

## Usage
1. Install dependencies:
    ```
    go mod tidy
    ```

2. Run the application:
    ```
    go run main.go
    ```

## Running Tests
To run the unit tests:

```
go test ./...
```