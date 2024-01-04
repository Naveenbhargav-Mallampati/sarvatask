package secondary

import (
	"context"
	"encoding/json"
	"filesarva/Application/ports"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"
)

func TestRaftConsensus(t *testing.T) {

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)
		var receivedData map[string]string
		err := json.NewDecoder(r.Body).Decode(&receivedData)
		assert.NoError(t, err)

		expectedFileName := "testFile.txt"
		expectedFileSize := "123"
		assert.Equal(t, expectedFileName, receivedData["key"])
		assert.Equal(t, expectedFileSize, receivedData["value"])
		w.WriteHeader(http.StatusOK)
	}))

	defer mockServer.Close()
	ctx := context.Background()

	testFileInfo := &ports.FileInfo{
		FileName: "testFile.txt",
		FileSize: 123,
	}
	logger := hclog.New(&hclog.LoggerOptions{
		Output: os.Stderr,
		Level:  hclog.Debug,
	})

	ctx = hclog.WithContext(ctx, logger)

	response, err := RaftConsensus(ctx, testFileInfo)
	assert.NoError(t, err)
	assert.Equal(t, "{\"message\":\"Dulicate:K/V\",\"data\":\"testFile.txt: 123\"}\n", response)
}
