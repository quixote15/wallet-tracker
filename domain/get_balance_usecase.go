package domain

import (
	"context"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// EthereumClient interface for dependency injection and testing
type EthereumClient interface {
	BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error)
}

type GetBalanceUseCase struct {
	client EthereumClient
}

type WalletBalance struct {
	Address string  `json:"address"`
	Balance float64 `json:"balance"`
}

func NewGetBalanceUseCase(client EthereumClient) *GetBalanceUseCase {
	return &GetBalanceUseCase{client: client}
}

func (uc *GetBalanceUseCase) Execute(ctx context.Context, address string) (*WalletBalance, error) {
	// Convert address string to Ethereum address
	account := common.HexToAddress(address)

	// Get balance using the passed context
	balance, err := uc.client.BalanceAt(ctx, account, nil)
	if err != nil {
		return nil, err
	}

	// Convert balance from Wei to Ether
	fBalance := new(big.Float)
	fBalance.SetString(balance.String())
	ethValue := new(big.Float).Quo(fBalance, big.NewFloat(math.Pow10(18)))

	floatBalance, _ := ethValue.Float64()

	return &WalletBalance{
		Address: address,
		Balance: floatBalance,
	}, nil
}
