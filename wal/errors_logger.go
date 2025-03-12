/*
*
*
**/

package wal

import (
	"log"
)

type ErrorsLogger struct {
	errLogger *log.Logger
}

func NewErrorsLogger() *ErrorsLogger {
	return &ErrorsLogger{
		errLogger: log.New("errors.log", "ERROR: ", log.Ldate|log.Ltime),
	}
}

func (e *ErrorsLogger) ReportError(err error) {
	e.errLogger.Println(err.Error())
}