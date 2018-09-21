package log

import (
	"fmt"
	"log"
	"os"
)

var (
	Info         *log.Logger
	Error        *log.Logger
	Warn         *log.Logger
	InfoHandler  = os.Stdout
	ErrorHandler = os.Stderr
	WarnHandler  = os.Stdout
	Namespace    = "zebedee-utils"
)

func init() {
	Info = log.New(InfoHandler, prefix("INFO"), log.Ldate|log.Ltime)
	Warn = log.New(WarnHandler, prefix("WARN"), log.Ldate|log.Ltime)
	Error = log.New(ErrorHandler, prefix("ERROR"), log.Ldate|log.Ltime|log.Lshortfile)
}

func prefix(level string) string {
	return fmt.Sprintf("[%s] %s: ", Namespace, level)
}
