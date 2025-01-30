package logger

// Константы для уровней логирования в строковом представлении
const (
	DebugStr = "DEBUG"
	InfoStr  = "INFO"
	WarnStr  = "WARN"
	ErrorStr = "ERROR"
	FatalStr = "FATAL"
)

// Константы для цветов в консоли
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorGray   = "\033[37m"
)

// Константы для форматирования
const (
	DefaultTimeFormat = "2006-01-02 15:04:05"
	DefaultSeparator  = " | "
	DefaultPrefix     = "APP"
)

// Константы для файлового логгера
const (
	DefaultMaxSize    = 100 // мегабайты
	DefaultMaxBackups = 3   // количество файлов
	DefaultMaxAge     = 28  // дни
	DefaultFilePerm   = 0644
)
