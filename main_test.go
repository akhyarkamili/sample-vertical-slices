package main

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServerHealthz(t *testing.T) {
	cfg := Config{
		Port: "13003",
	}

	go start(cfg)
	assert.Eventually(t, func() bool {
		resp, err := http.Get("http://localhost:13003/healthz")
		if err != nil {
			return false
		}
		return resp.StatusCode == http.StatusOK
	}, 1*time.Second, 200*time.Millisecond)
}

func TestServerPropose(t *testing.T) {
	cfg := Config{
		Port: "13004",
	}

	go start(cfg)
	require.Eventually(t, func() bool {
		resp, err := http.Get("http://localhost:13004/healthz")
		if err != nil {
			return false
		}
		return resp.StatusCode == http.StatusOK
	}, 200*time.Millisecond, 50*time.Millisecond)

	t.Run("propose endpoint exists and obeys contract", func(t *testing.T) {
		time.AfterFunc(200*time.Millisecond, func() {
			resp, err := http.Post("http://localhost:13004/propose", "application/json", strings.NewReader(`{"borrower_id": 1, "rate": 10, "principal_amount": 1000000}`))
			require.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			var response map[string]interface{}
			require.NoError(t, json.Unmarshal(body, &response))
			assert.Equal(t, true, response["success"])
			assert.NotEmpty(t, response["id"])
		})
	})
}
