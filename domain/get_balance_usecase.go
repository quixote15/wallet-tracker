package domain

import (
	"context"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type GetBalanceUseCase struct {
	client *ethclient.Client
}

type WalletBalance struct {
	Address string  `json:"address"`
	Balance float64 `json:"balance"`
}

func NewGetBalanceUseCase(client *ethclient.Client) *GetBalanceUseCase {
	return &GetBalanceUseCase{client: client}
}

func (uc *GetBalanceUseCase) Execute(address string) (*WalletBalance, error) {
	// Convert address string to Ethereum address
	account := common.HexToAddress(address)

	// Get balance
	balance, err := uc.client.BalanceAt(context.Background(), account, nil)
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