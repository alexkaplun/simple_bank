package server

import "go.uber.org/zap"

func Init() {

	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{"log/simplebank.log"}
	logger, _ := cfg.Build()

	logger.Info("Starting Server")

	r := NewRouter(logger)
	r.Run(":" +
		"8080")
	defer logger.Sync()
	defer logger.Info("Stopping Server")
}
