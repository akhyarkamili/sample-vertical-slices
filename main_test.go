package main

import (
	"encoding/json"
	"io"
	"loan-management/testhelper"
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

	go start(cfg, nil)
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
	db := testhelper.SetupDB(t)

	go start(cfg, db)
	time.Sleep(200 * time.Millisecond)

	t.Run("propose endpoint exists and obeys contract", func(t *testing.T) {
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
}

func TestServerApprove(t *testing.T) {
	cfg := Config{
		Port: "13005",
	}
	db := testhelper.SetupDB(t)

	go start(cfg, db)
	time.Sleep(200 * time.Millisecond)

	t.Run("approve endpoint exists and obeys contract", func(t *testing.T) {
		// Arrange
		proposeResp, err := http.Post("http://localhost:13005/propose", "application/json", strings.NewReader(`{"borrower_id": 1, "rate": 10, "principal_amount": 1000000}`))
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, proposeResp.StatusCode)
		proposeBody, err := io.ReadAll(proposeResp.Body)
		require.NoError(t, err)
		jsonResult := map[string]interface{}{}
		require.NoError(t, json.Unmarshal(proposeBody, &jsonResult))
		id := jsonResult["id"].(string)

		// Act
		s := `{"id": "` + strings.TrimSpace(id) + `", "employee_id": 1, "proof": "https://google.com"}`
		resp, err := http.Post("http://localhost:13005/approve", "application/json", strings.NewReader(s))
		require.NoError(t, err)
		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		// Assert
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		var response map[string]interface{}
		require.NoError(t, json.Unmarshal(body, &response))
		assert.Equal(t, true, response["success"])
	})
}
