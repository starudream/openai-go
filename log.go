package openai

import (
	"fmt"
	"log"
	"os"

	"github.com/go-resty/resty/v2"
)

var (
	_logger *log.Logger

	_skip = 7
)

func init() {
	flag := log.LstdFlags
	if Debug() {
		flag |= log.Lshortfile
	}
	_logger = log.New(os.Stdout, "", flag)
}

type restyLogger struct {
}

var _ resty.Logger = (*restyLogger)(nil)

func (l *restyLogger) Errorf(format string, v ...any) {
	l.output("ERR [HTTP] "+format, v...)
}

func (l *restyLogger) Warnf(format string, v ...any) {
	l.output("WAN [HTTP] "+format, v...)
}

func (l *restyLogger) Debugf(format string, v ...any) {
	l.output("DBG [HTTP] "+format, v...)
}

func (l *restyLogger) output(format string, v ...any) {
	if len(v) == 0 {
		_ = _logger.Output(_skip, format)
		return
	}
	_ = _logger.Output(_skip, fmt.Sprintf(format, v...))
}
