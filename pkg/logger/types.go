package customLogger

import (
	"log"
	"os"
)

type Error string

func (e Error) Error() string {
	return string(e)
}

var (
	ErrFileWriting = Error("ошибка записи в лог")
	ErrFileOpening = Error("ошибка открытия лог файла")
	ErrFileClosing = Error("ошибка закрытия лог файла")
)

const (
	permission = 0o600
	chunkSize  = 100
)

type Logger struct {
	logger  *log.Logger
	file    *os.File
	logChan chan string
}
