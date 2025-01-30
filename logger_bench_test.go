package logger

import (
	"testing"
)

func BenchmarkLogger(b *testing.B) {
	log := New().WithConfig(DefaultConfig())
	defer func() {
		if err := log.Close(); err != nil {
			b.Errorf("failed to close logger: %v", err)
		}
	}()

	b.Run("простое логирование", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			log.Info("тестовое сообщение")
		}
	})

	b.Run("логирование с полями", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			log.WithFields(map[string]interface{}{
				"key1": "value1",
				"key2": 123,
			}).Info("тестовое сообщение")
		}
	})

	b.Run("логирование с префиксом", func(b *testing.B) {
		prefixedLog := log.WithPrefix("TEST")
		for i := 0; i < b.N; i++ {
			prefixedLog.Info("тестовое сообщение")
		}
	})

	b.Run("форматированное логирование", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			log.Infof("сообщение %d: %s", i, "тест")
		}
	})
}
