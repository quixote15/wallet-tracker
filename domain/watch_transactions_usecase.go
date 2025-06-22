package domain

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Transaction struct {
	Hash        string    `json:"hash"`
	From        string    `json:"from"`
	To          string    `json:"to"`
	Value       string    `json:"value"`
	GasPrice    string    `json:"gas_price"`
	GasUsed     uint64    `json:"gas_used"`
	BlockNumber uint64    `json:"block_number"`
	Timestamp   time.Time `json:"timestamp"`
}

type WatchTransactionsUseCase struct {
	client *ethclient.Client
}

func NewWatchTransactionsUseCase(client *ethclient.Client) *WatchTransactionsUseCase {
	return &WatchTransactionsUseCase{client: client}
}

func (uc *WatchTransactionsUseCase) Execute(ctx context.Context, address string) error {
	// Convert address string to Ethereum address
	account := common.HexToAddress(address)

	// Get current block number as starting point
	startBlock, err := uc.client.BlockNumber(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get current block number: %v", err)
	}

	txChannel := make(chan Transaction)

	var wg sync.WaitGroup

	wg.Add(1)
	go uc.monitorTransactions(account, startBlock, txChannel, &wg)

	go func() {
		for transaction := range txChannel {
			fmt.Printf("Transaction from %s to %s\n", transaction.From, transaction.To)
		}
	}()

	wg.Wait()
	close(txChannel)

	return nil
}

func (uc *WatchTransactionsUseCase) monitorTransactions(account common.Address, startBlock uint64, txChannel chan Transaction, wg *sync.WaitGroup) {
	currentBlock := startBlock

	// Monitor for new blocks for a specified duration (5 minutes)
	monitorDuration := 5 * time.Minute
	startTime := time.Now()
	ticker := time.NewTicker(5 * time.Second) // Check for new blocks every 5 seconds
	defer ticker.Stop()

	for time.Since(startTime) < monitorDuration {
		select {
		case <-ticker.C:
			// Check for new blocks
			latestBlock, err := uc.client.BlockNumber(context.Background())
			if err != nil {
				continue
			}

			uc.processBlocks(account, currentBlock+1, latestBlock, txChannel)

			// Update current block to latest processed
			currentBlock = latestBlock
		default:
			// Small sleep to prevent busy waiting
			time.Sleep(100 * time.Millisecond)
		}
	}

	wg.Done()

}

func (uc *WatchTransactionsUseCase) processBlocks(account common.Address, startBlock, endBlock uint64, txChan chan Transaction) {

	for blockNum := startBlock; blockNum <= endBlock; blockNum++ {
		block, err := uc.client.BlockByNumber(context.Background(), big.NewInt(int64(blockNum)))
		if err != nil {
			continue // Skip blocks that can't be retrieved
		}

		for _, tx := range block.Transactions() {
			if transaction, isRelevant := uc.processTransaction(tx, account, blockNum, block.Time()); isRelevant {
				txChan <- transaction
			}
		}
	}

}

func (uc *WatchTransactionsUseCase) processTransaction(tx *types.Transaction, account common.Address, blockNum uint64, blockTime uint64) (Transaction, bool) {
	// Check if transaction involves our address (both from and to)
	var isRelevant bool
	if tx.To() != nil && tx.To().Hex() == account.Hex() {
		isRelevant = true
	}

	// Also check if the address is the sender
	signer := types.LatestSignerForChainID(tx.ChainId())
	from, err := types.Sender(signer, tx)
	if err == nil && from.Hex() == account.Hex() {
		isRelevant = true
	}

	if !isRelevant {
		return Transaction{}, false
	}

	// Get transaction receipt for gas used
	receipt, err := uc.client.TransactionReceipt(context.Background(), tx.Hash())
	if err != nil {
		return Transaction{}, false
	}

	toAddress := ""
	if tx.To() != nil {
		toAddress = tx.To().Hex()
	}

	transaction := Transaction{
		Hash:        tx.Hash().Hex(),
		From:        from.Hex(),
		To:          toAddress,
		Value:       tx.Value().String(),
		GasPrice:    tx.GasPrice().String(),
		GasUsed:     receipt.GasUsed,
		BlockNumber: blockNum,
		Timestamp:   time.Unix(int64(blockTime), 0),
	}

	return transaction, true
}
