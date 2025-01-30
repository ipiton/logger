package logger

import (
	"sync"
	"testing"

	"strings"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoggerConcurrentUsage(t *testing.T) {
	logger := New().
		WithLevel("debug").
		WithFile("stress.log")
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(num int) {
			defer wg.Done()
			logger.WithFields(map[string]interface{}{"goroutine": num}).
				Infof("Сообщение из горутины %d", num)
		}(i)
	}

	wg.Wait()

	require.NoError(t, logger.(*Logger).Sync(), "Ошибка синхронизации файла")

	content := readLogFileSecure(t, "stress.log")
	lines := strings.Count(content, "\n")
	assert.GreaterOrEqual(t, lines, 100, "Должно быть записано минимум 100 сообщений")
}
