package customLogger

import (
	"ads-service/pkg/utils"
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

// Init - Инициализация логгера.
func NewLogger(ctx context.Context) (Logger, error) {
	logFile := fmt.Sprintf("storage/logs/%s.log", time.Now().Format("2006-01-02"))
	var wg sync.WaitGroup
	logCh := make(chan string, chunkSize)

	if !utils.IsSafeLogPath(logFile) {
		return Logger{}, ErrFileOpening
	}

	// #nosec G304 -- путь проверяется функцией IsSafeLogPath
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_RDWR, permission)
	if err != nil {
		return Logger{}, ErrFileOpening
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer func() {
			if cerr := file.Close(); cerr != nil {
				log.Println("Error file closing") // log ERROR
			}
		}()

		for {
			select {
			case msg := <-logCh:
				if _, err := file.WriteString(msg); err != nil {
					log.Printf("Ошибка записи в лог: %v\n", err)
				}
			case <-ctx.Done():
				for {
					select {
					case msg := <-logCh:
						_, err := file.WriteString(msg)
						if err != nil {
							return
						}
					default:
						return
					}
				}
			}
		}
	}()

	return Logger{
		logger:  log.New(os.Stdout, "", log.LstdFlags|log.Ltime),
		file:    file,
		logChan: logCh,
	}, nil
}

func (l *Logger) ERROR(v ...interface{}) {
	currentTime := time.Now().Format(time.RFC3339)
	msg := fmt.Sprint(v...)
	log.Println(currentTime, msg)
	select {
	case l.logChan <- fmt.Sprintf("[%v] [ERROR]: %v\n", currentTime, msg):
	default:
		log.Printf("[%v] [ERROR]: logging channel overflow", currentTime)
	}
}

func (l *Logger) INFO(v ...interface{}) {
	currentTime := time.Now().Format(time.RFC3339)
	msg := fmt.Sprint(v...)
	select {
	case l.logChan <- fmt.Sprintf("[%v] [INFO]: %v\n", currentTime, msg):
	default:
		log.Printf("[%v] [INFO]: logging channel overflow", currentTime)
	}
}

func (l *Logger) WARN(v ...interface{}) {
	currentTime := time.Now().Format(time.RFC3339)
	msg := fmt.Sprint(v...)
	select {
	case l.logChan <- fmt.Sprintf("[%v] [WARN]: %v\n", currentTime, msg):
	default:
		log.Printf("[%v] [WARN]: logging channel overflow", currentTime)
	}
}
