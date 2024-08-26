package impl_test

import (
	response "balance_kafka_service/internal/domain/balance_service_response"
	"balance_kafka_service/internal/usecase/impl"
	"balance_kafka_service/internal/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetBalanceByID(t *testing.T) {
	t.Run("Success case", func(t *testing.T) {
		expectedResponse := response.GlobalResponse{
			StatusCode: http.StatusOK,
			Data:       map[string]interface{}{"balance": 1000},
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(expectedResponse)
		}))
		defer server.Close()

		// Overwrite the URL for the test
		oldURL := fmt.Sprintf("http://localhost:8081/balance/%d", 1)
		defer func() { fmt.Sprintf(oldURL) }()

		uc := impl.NewBalanceRequestUseCase()
		success, statusCode, err := uc.GetBalanceByID(1)

		require.NoError(t, err)
		assert.True(t, success)
		assert.Equal(t, http.StatusOK, statusCode)
	})

	t.Run("Max retries reached", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusServiceUnavailable)
		}))
		defer server.Close()

		// Overwrite the URL for the test
		oldURL := fmt.Sprintf("http://localhost:8081/balance/%d", 1)
		defer func() { fmt.Sprintf(oldURL) }()

		uc := impl.NewBalanceRequestUseCase()
		success, statusCode, err := uc.GetBalanceByID(1)

		require.Error(t, err)
		assert.False(t, success)
		assert.Equal(t, http.StatusServiceUnavailable, statusCode)
		assert.Equal(t, utils.ErrMaxRetryReached, err)
	})

	t.Run("Non-retryable error", func(t *testing.T) {
		expectedResponse := response.GlobalResponse{
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(expectedResponse)
		}))
		defer server.Close()

		// Overwrite the URL for the test
		oldURL := fmt.Sprintf("http://localhost:8081/balance/%d", 1)
		defer func() { fmt.Sprintf(oldURL) }()

		uc := impl.NewBalanceRequestUseCase()
		success, statusCode, err := uc.GetBalanceByID(1)

		require.Error(t, err)
		assert.False(t, success)
		assert.Equal(t, http.StatusBadRequest, statusCode)
		assert.Equal(t, utils.ErrHttpNotRetryAble, err)
	})

	t.Run("Decode error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("invalid json"))
		}))
		defer server.Close()

		// Overwrite the URL for the test
		oldURL := fmt.Sprintf("http://localhost:8081/balance/%d", 1)
		defer func() { fmt.Sprintf(oldURL) }()

		uc := impl.NewBalanceRequestUseCase()
		success, statusCode, err := uc.GetBalanceByID(1)

		require.Error(t, err)
		assert.False(t, success)
		assert.Equal(t, http.StatusInternalServerError, statusCode)
		assert.Equal(t, utils.ErrJsonDecode, err)
	})

	t.Run("HTTP request error", func(t *testing.T) {
		uc := impl.NewBalanceRequestUseCase()
		success, statusCode, err := uc.GetBalanceByID(999)

		require.Error(t, err)
		assert.False(t, success)
		assert.Equal(t, http.StatusInternalServerError, statusCode)
		assert.Equal(t, utils.ErrHttpRequest, err)
	})
}
