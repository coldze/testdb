package handlers

import (
	"fmt"

	"github.com/coldze/testdb/logs"
)

//This loggerFactory can generate unique logger for each request, thus they can be traced in logs.
type LoggerFactory func(prefix string) logs.Logger

func NewDefaultLoggerFactory(defaultLogger logs.Logger) LoggerFactory {
	return func(id string) logs.Logger {
		return logs.NewPrefixedLogger(defaultLogger, fmt.Sprintf(" [%s] ", id))
	}
}
