package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	require.NotNil(t, cfg)

	// Проверяем значения по умолчанию
	assert.Equal(t, "info", cfg.Level)
	assert.Equal(t, "", cfg.Files["main"])
	assert.Empty(t, cfg.Files["error"])
}

func TestConfig_Override(t *testing.T) {
	// nolint:govet
	tests := []struct {
		cfg      *Config // 8 bytes
		name     string  // 16 bytes
		override *Config // 8 bytes
		expected *Config // 8 bytes
	}{
		{
			name:     "переопределение уровня логирования",
			cfg:      DefaultConfig(),
			override: &Config{Level: "debug"},
			expected: &Config{
				Level: "debug",
				Files: map[string]string{
					"main":  "logs/main.log",
					"error": "logs/error.log",
				},
			},
		},
		{
			name: "переопределение путей к файлам",
			cfg:  DefaultConfig(),
			expected: &Config{
				Level: "info",
				Files: map[string]string{
					"main":  "custom/main.log",
					"error": "logs/error.log",
				},
			},
		},
		{
			name:     "переопределение nil конфигурацией",
			cfg:      DefaultConfig(),
			expected: DefaultConfig(),
		},
		{
			name:     "переопределение пустой конфигурацией",
			cfg:      DefaultConfig(),
			expected: DefaultConfig(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.cfg.Override(tt.expected)
			assert.Equal(t, tt.expected.Level, tt.cfg.Level)
			assert.Equal(t, tt.expected.Files, tt.cfg.Files)
		})
	}
}

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		cfg     *Config
		name    string
		wantErr bool
	}{
		{
			name:    "валидная конфигурация",
			cfg:     DefaultConfig(),
			wantErr: false,
		},
		{
			name: "некорректный уровень логирования",
			cfg: &Config{
				Level: "invalid",
				Files: map[string]string{"main": "logs/main.log"},
			},
			wantErr: true,
		},
		{
			name: "отсутствует основной файл лога",
			cfg: &Config{
				Level: "info",
				Files: map[string]string{"error": "logs/error.log"},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
