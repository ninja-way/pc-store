package config

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

const LogDir = ".logs"

func SetupLogger(cfg *Config) {
	if cfg.Environment == "prod" {
		if err := os.MkdirAll(LogDir, os.ModePerm); err != nil {
			log.Fatal(err)
		}

		date := time.Now().Format("02-01-2006")
		logFile := fmt.Sprintf("%s/%s.log", LogDir, date)

		file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
		if err != nil {
			log.Fatal("failed open log file")
		}

		log.SetOutput(file)
		log.SetFormatter(&log.JSONFormatter{})
		log.SetLevel(log.WarnLevel)
		return
	}

	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}
