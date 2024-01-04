package secondary

import (
	"bytes"
	"context"
	"encoding/json"
	"filesarva/Application/ports"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/hashicorp/go-hclog"
)

func RaftConsensus(ctx context.Context, fileInfo *ports.FileInfo) (string, error) {
	logger := hclog.FromContext(ctx)
	url := "http://127.0.0.1:9090"
	fmt.Println(fileInfo)
	data1 := struct {
		Message string `json:"key"`
		Data    string `json:"value"`
	}{
		Message: fileInfo.FileName,
		Data:    strconv.FormatInt(fileInfo.FileSize, 10),
	}
	payload, err := json.Marshal(data1)
	logger.Info("Added payload")
	if err != nil {
		errorMessage := fmt.Sprintf("error encoding JSON payload: %v", err)
		logger.Error(errorMessage)
		return "", fmt.Errorf("error encoding JSON payload: %v", err)
	}

	request, err := http.NewRequest("PUT", url, bytes.NewBuffer(payload))
	if err != nil {
		errorMessage := fmt.Sprintf("error creating request:  %v", err)
		logger.Error(errorMessage)
		return "", fmt.Errorf("error creating request: %v", err)
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		errorMessage := fmt.Sprintf("error making request:  %v", err)
		logger.Error(errorMessage)
		return "", fmt.Errorf("error making request: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		errorMessage := fmt.Sprintf("Unexpected Status Code returned:  %s", response.Status)
		logger.Error(errorMessage)
		return "", fmt.Errorf("unexpected status code: %s", response.Status)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		errorMessage := fmt.Sprintf("error reading response body:  %v", err)
		logger.Error(errorMessage)
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	responseString := string(body)
	logger.Debug("Raft consensus completed", "filename", fileInfo.FileName, "size", fileInfo.FileSize, "response", responseString)
	return responseString, nil
}
