package main

import (
	"github.com/BurntSushi/toml"
	"github.com/sirupsen/logrus"
	"os"
	"src/configs/backend"
	_ "src/docs"
	"src/server"
)

var configPath = os.Getenv("CONFIG_FILE")

// @title           BMSTU-WEB API
// @version         1.0
// @description     Server for bmstu dormitory app.
// @host      localhost:8082
func main() {
	config := configs.CreateConfigForServer()
	_, err := toml.DecodeFile(configPath, &config)
	if err != nil {
		logrus.Fatal(err)
	}
	contextLogger := logrus.WithFields(logrus.Fields{})
	logrus.SetReportCaller(false)
	logrus.SetFormatter(&logrus.TextFormatter{PadLevelText: false, DisableLevelTruncation: false})
	appServer := server.CreateServer(config, contextLogger)

	err = appServer.Start()
	if err != nil {
		panic(err)
	}
}
