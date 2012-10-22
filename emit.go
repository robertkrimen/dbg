package dbg

import (
	"fmt"
	"io"
	"log"
	"os"
)

func ln(tmp string) string {
	length := len(tmp)
	if length > 0 && tmp[length-1] != '\n' {
		return tmp + "\n"
	}
	return tmp
}

type _emit interface {
	emit(_frmt, string, ...interface{})
}

type _emitWriter struct {
	writer io.Writer
}

func (self _emitWriter) emit(frmt _frmt, format string, values ...interface{}) {
	if format == "" {
		fmt.Fprintln(self.writer, values...)
	} else {
		if frmt.panic {
			panic(fmt.Sprintf(format, values...))
		}
		fmt.Fprintf(self.writer, ln(format), values...)
		if frmt.fatal {
			os.Exit(1)
		}
	}
}

type _emitLogger struct {
	logger *log.Logger
}

func (self _emitLogger) emit(frmt _frmt, format string, values ...interface{}) {
	if format == "" {
		self.logger.Println(values...)
	} else {
		if frmt.panic {
			self.logger.Panicf(format, values...)
		} else if frmt.fatal {
			self.logger.Fatalf(format, values...)
		} else {
			self.logger.Printf(format, values...)
		}
	}
}

type _emitLog struct {
}

func (self _emitLog) emit(frmt _frmt, format string, values ...interface{}) {
	if format == "" {
		log.Println(values...)
	} else {
		if frmt.panic {
			log.Panicf(format, values...)
		} else if frmt.fatal {
			log.Fatalf(format, values...)
		} else {
			log.Printf(format, values...)
		}
	}
}
