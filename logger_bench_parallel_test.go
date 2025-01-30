package logger

import (
	"testing"
)

func BenchmarkLoggerParallel(b *testing.B) {
	log := New().WithConfig(DefaultConfig())
	defer func() {
		if err := log.Close(); err != nil {
			b.Errorf("failed to close logger: %v", err)
		}
	}()

	b.Run("параллельное логирование", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				log.Info("тестовое сообщение")
			}
		})
	})

	b.Run("параллельное логирование с полями", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				log.WithFields(map[string]interface{}{
					"goroutine": "test",
					"count":     b.N,
				}).Info("тестовое сообщение")
			}
		})
	})

	b.Run("параллельное логирование с префиксом", func(b *testing.B) {
		prefixedLog := log.WithPrefix("TEST")
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				prefixedLog.Info("тестовое сообщение")
			}
		})
	})
}
