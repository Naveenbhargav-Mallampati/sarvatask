package main

import (
	"filesarva/Application/adapters/primary/api"

	"github.com/hashicorp/go-hclog"
	"gopkg.in/natefinch/lumberjack.v2"
)

func main() {
	logger := hclog.New(&hclog.LoggerOptions{
		Name:  "file-upload-service",
		Level: hclog.Debug,
		// Use lumberjack as the log handler
		Output: &lumberjack.Logger{
			Filename:   "./Debug/file-upload-service.log",
			MaxSize:    50, // megabytes
			MaxBackups: 3,
			MaxAge:     10, // days
			Compress:   true,
		},
	})
	api.Driver(logger)

}
