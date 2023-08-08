package helper

import (
	"fmt"
	"io"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

func InitLogger() {
	var err error
	if len(AppConfig.Logging.LogFile) < 3 {
		AppConfig.Logging.LogFile, err = os.Getwd()
		if err != nil {
			fmt.Printf("Runtime path error: %v", err)
			os.Exit(1)
		}
		AppConfig.Logging.LogFile += "/vmausers.log"
	}

	lumberjackLogger := &lumberjack.Logger{
		Filename:   AppConfig.Logging.LogFile,
		MaxSize:    AppConfig.Logging.MaxSizeMB,
		MaxBackups: AppConfig.Logging.MaxBackups,
		MaxAge:     AppConfig.Logging.MaxAge,
	}

	// Fork writing into two outputs
	multiWriter := io.MultiWriter(os.Stderr, lumberjackLogger)

	logFormatter := new(log.TextFormatter)
	logFormatter.TimestampFormat = time.RFC1123Z
	logFormatter.FullTimestamp = true

	log.SetFormatter(logFormatter)
	log.SetLevel(log.InfoLevel)
	log.SetOutput(multiWriter)
}
