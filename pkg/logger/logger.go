package customLogger

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"
)

type Logger struct {
	file   *os.File
	logger *log.Logger
}

// NewLogger - Инициализация логгера.
func NewLogger() (Logger, error) {
	f, err := os.OpenFile("storage/ad-service.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666) // filename string
	if err != nil {
		log.Println(err)
		return Logger{}, errors.New("Failed to open log file")
	}
	return Logger{
		logger: log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Llongfile),
		file:   f,
	}, nil
}

func (l *Logger) ERROR(v ...interface{}) {
	currentTime := time.Now().Format(time.RFC3339)
	msg := fmt.Sprint(v...)
	_, err := l.file.WriteString(currentTime + " [ERROR] " + msg + "\n")
	if err != nil {
		log.Println(err)
	}
}

func (l *Logger) INFO(v ...interface{}) {
	currentTime := time.Now().Format(time.RFC3339)
	msg := fmt.Sprint(v...)
	_, err := l.file.WriteString(currentTime + " [INFO] " + msg + "\n")
	if err != nil {
		log.Println(err)
	}
}

func (l *Logger) WARN(v ...interface{}) {
	currentTime := time.Now().Format(time.RFC3339)
	msg := fmt.Sprint(v...)
	_, err := l.file.WriteString(currentTime + " [WARN] " + msg + "\n")
	if err != nil {
		log.Println(err)
	}
}
