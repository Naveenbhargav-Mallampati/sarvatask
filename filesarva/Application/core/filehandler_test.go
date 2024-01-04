// package core

// import (
// 	"context"
// 	"filesarva/Application/ports"
// 	"fmt"
// 	"os"
// 	"path/filepath"
// 	"testing"

// 	"github.com/hashicorp/go-hclog"
// )

// type mockRaftConsensus struct{}

// func (m *mockRaftConsensus) RaftConsensus(ctx context.Context, fileInfo *ports.FileInfo) (string, error) {
// 	// Mock implementation of RaftConsensus for testing
// 	return "mocked_output", nil
// }

// func TestProcessFile(t *testing.T) {

// 	t.Run("ProcessNonEmptyFile", func(t *testing.T) {
// 		// Get the absolute path to the Schedule.md in the filesarva directory
// 		wd, err := os.Getwd()
// 		if err != nil {
// 			t.Fatalf("Error getting working directory: %v", err)
// 		}
// 		fmt.Println("Current working directory:", wd)

// 		logger := hclog.New(nil)
// 		ctx := hclog.WithContext(context.Background(), logger)
// 		filePath := filepath.Join("C:/Users/vasav/Documents/GitHub/filesarva/Schedule.md")

// 		// Run the function being tested
// 		ProcessFile(ctx, filePath)
// 		// Add assertions here based on the expected behavior
// 		// You might want to use the testing package's functions like t.Errorf, etc.
// 	})
// 	t.Run("ProcessEmptyFile", func(t *testing.T) {

// 		logger := hclog.New(nil)
// 		ctx := hclog.WithContext(context.Background(), logger)
// 		filePath := filepath.Join("C:/Users/vasav/Documents/GitHub/filesarva/empty.txt")

//			// Run the function being tested
//			ProcessFile(ctx, filePath)
//		})
//	}
package core

import (
	"context"
	"filesarva/Application/ports"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
)

type mockRaftConsensus struct{}

func (m *mockRaftConsensus) RaftConsensus(ctx context.Context, fileInfo *ports.FileInfo) (string, error) {

	return "mocked_output", nil
}

func TestProcessFile(t *testing.T) {

	t.Run("ProcessNonEmptyFile", func(t *testing.T) {

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		ctx := echo.New().NewContext(req, rec)

		filePath := filepath.Join("C:/Users/vasav/Documents/GitHub/filesarva/preview.jpg")
		output, err := ProcessFile(ctx, filePath)
		output = strings.TrimSpace(output)
		if err != nil {
			t.Errorf("ProcessFile returned an error: %v", err)
		}

		expectedOutput := `{"message":"Dulicate:K/V","data":"preview.jpg: 297708"}`
		if output != expectedOutput {
			t.Errorf("Unexpected output. Got %s expected %s", output, expectedOutput)
		}

	})

	t.Run("ProcessEmptyFile", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		ctx := echo.New().NewContext(req, rec)

		filePath := filepath.Join("C:/Users/vasav/Documents/GitHub/filesarva/empty.txt")

		output, err := ProcessFile(ctx, filePath)

		if err == nil {
			t.Error("Expected an error for an empty file, but got nil")
		}

		if output != "" {
			t.Errorf("Expected empty output for an empty file, but got %s", output)
		}
	})
}
