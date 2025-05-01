/*
*
*
**/

package wal

import (
	"log"
	"os"
	"fmt"
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
	fmt.Println(err)
	e.errLogger.Println(err.Error())
}