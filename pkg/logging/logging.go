package logging

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorGrey   = "\033[90m"
)

type LogMessage struct {
	Time     time.Time
	Tick     int
	Color    string
	Level    string
	Filename string
	Message  string
}

type logger struct {
	queue    chan LogMessage
	filepath string
	tick     int
	wg       sync.WaitGroup
}

func (l *logger) run() {
	log.SetOutput(
		&lumberjack.Logger{
			Filename:   l.filepath,
			MaxSize:    30,
			MaxBackups: 3,
			MaxAge:     14,
			Compress:   true,
		},
	)
	log.SetFlags(0)

	for msg := range l.queue {
		timeTick := fmt.Sprintf("%s%s", colorGrey, msg.Time.Format("15:04:05.000"))
		lines := strings.SplitSeq(msg.Message, "\n")
		for line := range lines {
			log.Printf("%s %d %s[%s] %s %s%s\n", timeTick, l.tick, msg.Color, msg.Level, msg.Filename, line, colorReset)
		}
	}
	l.wg.Done()
}

func Init(filepath string, tick int) {
	l := &logger{
		queue:    make(chan LogMessage, 10),
		filepath: filepath,
		tick:     tick,
	}

	l.wg.Add(1)
	go l.run()
	_logger = l
}

func Flush() {
	close(_logger.queue)
	_logger.wg.Wait()
}

func (l *logger) log(level, color, format string, args ...any) {
	_, file, _, ok := runtime.Caller(2)
	fileName := "UNKNOWN"
	if ok {
		fileName = strings.ToUpper(filepath.Base(file))
	}

	var msg string
	if len(args) > 0 {
		msg = fmt.Sprintf(format, args...)
	} else {
		msg = format
	}

	l.queue <- LogMessage{
		Time:     time.Now(),
		Tick:     l.tick,
		Level:    level,
		Filename: fileName,
		Message:  msg,
		Color:    color,
	}
}

var _logger *logger

func Info(format string, args ...any)  { _logger.log("TELEMETRY", colorGrey, format, args...) }
func Debug(format string, args ...any) { _logger.log("TELEMETRY", colorReset, format, args...) }
func Error(format string, args ...any) { _logger.log("FAULT", colorRed, format, args...) }
func Warn(format string, args ...any)  { _logger.log("WARN", colorYellow, format, args...) }
func Ok(format string, args ...any)    { _logger.log("STABLE", colorGreen, format, args...) }
