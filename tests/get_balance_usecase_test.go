package tests

import (
	"context"
	"math/big"
	"testing"
	"wallet-tracker/domain"

	"github.com/ethereum/go-ethereum/common" // Make sure to run: go get github.com/ethereum/go-ethereum
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockEthClient implements domain.EthereumClient interface using testify/mock
type MockEthClient struct {
	mock.Mock
}

func (m *MockEthClient) BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error) {
	args := m.Called(ctx, account, blockNumber)
	return args.Get(0).(*big.Int), args.Error(1)
}

func TestGetBalanceUseCase(t *testing.T) {
	t.Run("Should get balance from an address and verify method call", func(t *testing.T) {
		// Arrange
		mockedClient := new(MockEthClient)
		useCase := domain.NewGetBalanceUseCase(mockedClient)
		address := "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6"
		expectedAccount := common.HexToAddress(address)
		
		// Setup mock expectation - Return 1 ETH in Wei
		oneEthInWei := new(big.Int)
		oneEthInWei.SetString("1000000000000000000", 10)
		mockedClient.On("BalanceAt", mock.Anything, expectedAccount, (*big.Int)(nil)).Return(oneEthInWei, nil)

		// Act
		result, err := useCase.Execute(context.Background(), address)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, 1.0, result.Balance)
		assert.Equal(t, address, result.Address)
		
		// Verify that BalanceAt was called exactly once with expected parameters
		mockedClient.AssertExpectations(t)
		mockedClient.AssertCalled(t, "BalanceAt", mock.Anything, expectedAccount, (*big.Int)(nil))
	})

	t.Run("Should return error when BalanceAt fails", func(t *testing.T) {
		// Arrange
		mockedClient := new(MockEthClient)
		useCase := domain.NewGetBalanceUseCase(mockedClient)
		address := "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6"
		expectedAccount := common.HexToAddress(address)
		
		// Setup mock to return error
		mockedClient.On("BalanceAt", mock.Anything, expectedAccount, (*big.Int)(nil)).Return((*big.Int)(nil), assert.AnError)

		// Act
		result, err := useCase.Execute(context.Background(), address)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, assert.AnError, err)
		
		// Verify method was called
		mockedClient.AssertExpectations(t)
	})

	t.Run("Should handle different balance amounts", func(t *testing.T) {
		// Arrange
		mockedClient := new(MockEthClient)
		useCase := domain.NewGetBalanceUseCase(mockedClient)
		address := "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6"
		expectedAccount := common.HexToAddress(address)
		
		// Setup mock to return 0.5 ETH in Wei
		halfEthInWei := new(big.Int)
		halfEthInWei.SetString("500000000000000000", 10) // 0.5 * 10^18
		mockedClient.On("BalanceAt", mock.Anything, expectedAccount, (*big.Int)(nil)).Return(halfEthInWei, nil)

		// Act
		result, err := useCase.Execute(context.Background(), address)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, 0.5, result.Balance)
		assert.Equal(t, address, result.Address)
		
		// Verify method was called
		mockedClient.AssertExpectations(t)
	})

	t.Run("Should verify BalanceAt called with correct context", func(t *testing.T) {
		// Arrange
		mockedClient := new(MockEthClient)
		useCase := domain.NewGetBalanceUseCase(mockedClient)
		address := "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6"
		expectedAccount := common.HexToAddress(address)
		customCtx := context.WithValue(context.Background(), "test", "value")
		
		// Setup mock with specific context expectation
		oneEthInWei := new(big.Int)
		oneEthInWei.SetString("1000000000000000000", 10)
		mockedClient.On("BalanceAt", customCtx, expectedAccount, (*big.Int)(nil)).Return(oneEthInWei, nil)

		// Act
		result, err := useCase.Execute(customCtx, address)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, 1.0, result.Balance)
		
		// Verify method was called with the exact context
		mockedClient.AssertExpectations(t)
	})
}
