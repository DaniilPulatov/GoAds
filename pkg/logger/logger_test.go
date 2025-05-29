package customLogger

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLoggerMethods(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger, err := NewLogger(ctx)
	assert.NoError(t, err)

	// Проверяем, что методы не паникуют
	assert.NotPanics(t, func() { logger.ERROR("test error") })
	assert.NotPanics(t, func() { logger.INFO("test info") })
	assert.NotPanics(t, func() { logger.WARN("test warn") })

	time.Sleep(10 * time.Millisecond)
}

func TestLogger_ErrorOverflow(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger, err := NewLogger(ctx)
	assert.NoError(t, err)

	// Переполняем канал логов
	for i := range chunkSize + 10 {
		logger.ERROR("overflow test", i)
	}
	time.Sleep(10 * time.Millisecond)
}

func TestLogger_INFO(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger, err := NewLogger(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, logger)

	logger.INFO("test info message")

	// Проверяем, что логгер не nil и канал логов создан
	assert.NotNil(t, logger.logChan)
	assert.NotNil(t, logger.file)

	// Проверяем, что файл существует
	_, statErr := os.Stat(logger.file.Name())
	assert.NoError(t, statErr)
}
