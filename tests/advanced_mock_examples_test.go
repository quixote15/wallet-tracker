package tests

import (
	"context"
	"errors"
	"math/big"
	"testing"
	"time"
	"wallet-tracker/domain"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Advanced testing examples showing the power of testify/mock
func TestAdvancedMockingPatterns(t *testing.T) {
	t.Run("Should verify exact number of calls", func(t *testing.T) {
		// Arrange
		mockedClient := new(MockEthClient)
		useCase := domain.NewGetBalanceUseCase(mockedClient)
		address := "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6"
		expectedAccount := common.HexToAddress(address)

		// Setup mock to expect exactly 1 call
		oneEthInWei := new(big.Int)
		oneEthInWei.SetString("1000000000000000000", 10)
		mockedClient.On("BalanceAt", mock.Anything, expectedAccount, (*big.Int)(nil)).Return(oneEthInWei, nil).Once()

		// Act
		_, err := useCase.Execute(context.Background(), address)

		// Assert
		assert.NoError(t, err)
		mockedClient.AssertExpectations(t) // This will fail if called more or less than once
	})

	t.Run("Should handle multiple different return values", func(t *testing.T) {
		// Arrange
		mockedClient := new(MockEthClient)
		useCase := domain.NewGetBalanceUseCase(mockedClient)
		address1 := "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6"
		address2 := "0x123d35Cc6634C0532925a3b8D4C9db96C4b4d8b6"
		account1 := common.HexToAddress(address1)
		account2 := common.HexToAddress(address2)

		// Setup different return values for different addresses
		oneEthInWei := new(big.Int)
		oneEthInWei.SetString("1000000000000000000", 10)
		twoEthInWei := new(big.Int)
		twoEthInWei.SetString("2000000000000000000", 10)

		mockedClient.On("BalanceAt", mock.Anything, account1, (*big.Int)(nil)).Return(oneEthInWei, nil)
		mockedClient.On("BalanceAt", mock.Anything, account2, (*big.Int)(nil)).Return(twoEthInWei, nil)

		// Act
		result1, err1 := useCase.Execute(context.Background(), address1)
		result2, err2 := useCase.Execute(context.Background(), address2)

		// Assert
		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.Equal(t, 1.0, result1.Balance)
		assert.Equal(t, 2.0, result2.Balance)
		mockedClient.AssertExpectations(t)
	})

	t.Run("Should verify call order with multiple calls", func(t *testing.T) {
		// Arrange
		mockedClient := new(MockEthClient)
		useCase := domain.NewGetBalanceUseCase(mockedClient)
		address := "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6"
		expectedAccount := common.HexToAddress(address)

		// Setup mock to return different values on subsequent calls
		oneEthInWei := new(big.Int)
		oneEthInWei.SetString("1000000000000000000", 10)
		twoEthInWei := new(big.Int)
		twoEthInWei.SetString("2000000000000000000", 10)

		// First call returns 1 ETH, second call returns 2 ETH
		mockedClient.On("BalanceAt", mock.Anything, expectedAccount, (*big.Int)(nil)).Return(oneEthInWei, nil).Once()
		mockedClient.On("BalanceAt", mock.Anything, expectedAccount, (*big.Int)(nil)).Return(twoEthInWei, nil).Once()

		// Act
		result1, err1 := useCase.Execute(context.Background(), address)
		result2, err2 := useCase.Execute(context.Background(), address)

		// Assert
		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.Equal(t, 1.0, result1.Balance)
		assert.Equal(t, 2.0, result2.Balance)
		mockedClient.AssertExpectations(t)
	})

	t.Run("Should test with custom error types", func(t *testing.T) {
		// Arrange
		mockedClient := new(MockEthClient)
		useCase := domain.NewGetBalanceUseCase(mockedClient)
		address := "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6"
		expectedAccount := common.HexToAddress(address)
		customError := errors.New("custom network timeout error")

		// Setup mock to return custom error
		mockedClient.On("BalanceAt", mock.Anything, expectedAccount, (*big.Int)(nil)).Return((*big.Int)(nil), customError)

		// Act
		result, err := useCase.Execute(context.Background(), address)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, customError, err)
		assert.Contains(t, err.Error(), "custom network timeout")
		mockedClient.AssertExpectations(t)
	})

	t.Run("Should verify context timeout behavior", func(t *testing.T) {
		// Arrange
		mockedClient := new(MockEthClient)
		useCase := domain.NewGetBalanceUseCase(mockedClient)
		address := "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6"
		expectedAccount := common.HexToAddress(address)

		// Create a context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		// Setup mock to return context deadline exceeded
		mockedClient.On("BalanceAt", ctx, expectedAccount, (*big.Int)(nil)).Return((*big.Int)(nil), context.DeadlineExceeded)

		// Act
		result, err := useCase.Execute(ctx, address)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, context.DeadlineExceeded, err)
		mockedClient.AssertExpectations(t)
	})

	t.Run("Should use MatchedBy for complex argument matching", func(t *testing.T) {
		// Arrange
		mockedClient := new(MockEthClient)
		useCase := domain.NewGetBalanceUseCase(mockedClient)
		address := "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6"

		// Setup mock with custom matcher for context
		oneEthInWei := new(big.Int)
		oneEthInWei.SetString("1000000000000000000", 10)

		// Match any context that has a specific value
		contextMatcher := mock.MatchedBy(func(ctx context.Context) bool {
			return ctx != nil
		})

		// Match any address that starts with 0x742
		addressMatcher := mock.MatchedBy(func(addr common.Address) bool {
			return addr.Hex()[:5] == "0x742"
		})

		mockedClient.On("BalanceAt", contextMatcher, addressMatcher, (*big.Int)(nil)).Return(oneEthInWei, nil)

		// Act
		result, err := useCase.Execute(context.Background(), address)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, 1.0, result.Balance)
		mockedClient.AssertExpectations(t)
	})

	t.Run("Should verify method was never called", func(t *testing.T) {
		// Arrange
		mockedClient := new(MockEthClient)
		// Don't set up any expectations

		// Act - Don't call the method

		// Assert - Verify the method was never called
		mockedClient.AssertNotCalled(t, "BalanceAt")
	})
}