package controllers

import (
	"context"
	"encoding/json"
	"math"
	"math/big"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type WalletBalance struct {
	Address string  `json:"address"`
	Balance float64 `json:"balance"`
}

type WalletController struct {
	client *ethclient.Client
}

func NewWalletController(client *ethclient.Client) *WalletController {
	return &WalletController{client: client}
}

func (c *WalletController) Routes() []Route {
	const baseRoute = "/wallet"

	return []Route{
		{
			Path:    baseRoute + "/{address}/balance",
			Method:  http.MethodGet,
			Handler: c.GetBalance,
		},
	}
}

func (c *WalletController) GetBalance(w http.ResponseWriter, r *http.Request) {
	// Get address from URL parameters using chi's URLParam
	address := r.URL.Path[len("/wallet/"):len(r.URL.Path)-len("/balance")]

	// Convert address string to Ethereum address
	account := common.HexToAddress(address)

	// Get balance
	balance, err := c.client.BalanceAt(context.Background(), account, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert balance from wei to ether
	fbalance := new(big.Float)
	fbalance.SetString(balance.String())
	ethValue := new(big.Float).Quo(fbalance, big.NewFloat(math.Pow10(18)))

	// Convert to float64 for JSON marshaling
	ethFloat, _ := ethValue.Float64()

	// Prepare response
	response := WalletBalance{
		Address: address,
		Balance: ethFloat,
	}

	// Set content type and encode response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}