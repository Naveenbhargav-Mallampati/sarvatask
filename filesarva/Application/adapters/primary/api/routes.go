package api

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"filesarva/Application/core"

	"github.com/hashicorp/go-hclog"
	"github.com/labstack/echo/v4"
)

func Upload(c echo.Context) error {
	var wg sync.WaitGroup
	logger := hclog.FromContext(c.Request().Context())
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		errormsg := fmt.Sprintf("Failed to open file : %v", err)
		logger.Error(errormsg)
		return err
	}
	defer src.Close()

	dst, err := os.Create(file.Filename)
	if err != nil {
		errormsg := fmt.Sprintf("Failed to create file : %v", err)
		logger.Error(errormsg)
		return err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return err
	}
	logger.Info("Created file Copy")

	fileinfo, err := dst.Stat()
	if err != nil {

		return err
	}
	if fileinfo.Size() > 0 {
		logger.Info("The file is not empty.")
	} else {
		logger.Info("The file is empty.")
	}
	logger.Info("Checked Empty status")
	resultCh := make(chan string, 1)

	wg.Add(1)
	go func() {
		defer wg.Done()
		out, err := core.ProcessFile(c, fileinfo.Name())
		logger.Info("Sent to core")
		if err != nil {
			errormsg := fmt.Sprintf("ProcessFile failed: %v", err)
			logger.Error(errormsg)
			resultCh <- fmt.Sprintf("Error processing file: %s", err)
			return
		}
		resultCh <- out
	}()

	timeout := time.Second * 10
	select {
	case <-time.After(timeout):
		return c.HTML(http.StatusInternalServerError, "<p>Timeout processing file</p>")
	case output, ok := <-resultCh:
		if !ok {
			return c.HTML(http.StatusInternalServerError, "<p>Error processing file</p>")
		}

		wg.Wait()
		close(resultCh)
		logger.Info("returned the output")
		return c.HTML(http.StatusOK, fmt.Sprintf("<p>File %s uploaded successfully. RaftConsensus Output: %s</p>", file.Filename, output))
	}
}
