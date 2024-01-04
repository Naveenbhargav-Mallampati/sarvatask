package core

import (
	"filesarva/Application/adapters/secondary"
	"filesarva/Application/ports"
	"fmt"
	"os"
	"sync"

	"github.com/hashicorp/go-hclog"
	"github.com/labstack/echo/v4"
)

func ProcessFile(ctx echo.Context, filePath string) (string, error) {
	logger := hclog.FromContext(ctx.Request().Context())
	logger.Info("File processing started", "filename", filePath)
	var raftConsensusWG sync.WaitGroup
	file, err := os.Open(filePath)
	if err != nil {
		logger.Error("Error opening file: %s", err)
		return "", err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		logger.Error("Error getting file info: %s", err)
		return "", err
	}

	if fileInfo.Size() == 0 {
		logger.Error("File %s is empty", filePath)
		return "", fmt.Errorf("file is empty")
	}

	fileData := &ports.FileInfo{
		FileName: fileInfo.Name(),
		FileSize: fileInfo.Size(),
	}
	var output string
	outputCh := make(chan string, 1)

	raftConsensusWG.Add(1)
	go func() {
		defer raftConsensusWG.Done()
		httpRequest := ctx.Request()
		out, err := secondary.RaftConsensus(httpRequest.Context(), fileData)
		if err != nil {
			logger.Error("Raft consensus failed", "error", err)
		}
		output = out
		fmt.Println("In filehandler")
		fmt.Println(output)
		outputCh <- output
	}()

	raftConsensusWG.Wait()
	close(outputCh)
	logger.Debug("File processing completed", "filename", fileData.FileName, "size", fileData.FileSize)
	return output, nil
}
