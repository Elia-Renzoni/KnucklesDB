/*
*
*
**/

package wal

import (
	"log"
	"os"
)

type ErrorsLogger struct {
	errLogger *log.Logger
}

func NewErrorsLogger() *ErrorsLogger {
	return &ErrorsLogger{
		errLogger: log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime),
	}
}

func (e *ErrorsLogger) ReportError(err error) {
	e.errLogger.Println(err.Error())
}