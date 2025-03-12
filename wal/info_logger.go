package wal

import (
	"log"
	"os"
)

type InfoLogger struct {
	logger *log.Logger
}

func NewInfoLogger() *InfoLogger {
	return &InfoLogger{
		logger: log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime),
	}
}

func (i *InfoLogger) ReportInfo(info string) {
	i.logger.Println(info)
}
