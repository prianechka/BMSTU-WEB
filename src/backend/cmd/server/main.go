package main

import (
	"github.com/BurntSushi/toml"
	"github.com/sirupsen/logrus"
	"src/configs"
	_ "src/docs"
	"src/server"
)

var configPath = "/home/prianechka/Education/BMSTU/Web/BMSTU-WEB/src/backend/configs/config.toml"

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
