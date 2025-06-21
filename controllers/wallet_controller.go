package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/ethereum/go-ethereum/ethclient"

	"wallet-tracker/domain"
)

type WalletController struct {
	getBalanceUseCase *domain.GetBalanceUseCase
}

func NewWalletController(client *ethclient.Client) *WalletController {
	return &WalletController{
		getBalanceUseCase: domain.NewGetBalanceUseCase(client),
	}
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
	// Get address from URL parameters
	address := r.URL.Path[len("/wallet/") : len(r.URL.Path)-len("/balance")]

	// Get balance using the use case
	balance, err := c.getBalanceUseCase.Execute(address)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set content type and encode response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(balance)
}
