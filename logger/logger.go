package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

var (
	// Log is the global logger instance
	Log *logrus.Logger
	// telegramHook 全局Telegram Hook实例
	telegramHook *TelegramHook
)

// compactFormatter is a custom formatter for cleaner log output
type compactFormatter struct {
	logrus.TextFormatter
}

func (f *compactFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	level := strings.ToUpper(entry.Level.String())[0:4]

	// Skip frames to find actual caller (skip logrus + our wrapper functions)
	caller := ""
	for i := 3; i < 10; i++ {
		_, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		// Skip logrus internal and our logger.go
		if !strings.Contains(file, "logrus") && !strings.HasSuffix(file, "logger/logger.go") {
			// Get package name from path (e.g., "nofx/manager/trader_manager.go" -> "manager")
			dir := filepath.Dir(file)
			pkg := filepath.Base(dir)
			caller = fmt.Sprintf("%s/%s:%d", pkg, filepath.Base(file), line)
			break
		}
	}

	msg := fmt.Sprintf("[%s] %s %s\n", level, caller, entry.Message)
	return []byte(msg), nil
}

func init() {
	// Auto-initialize default logger to ensure it works before Init is called
	Log = logrus.New()
	Log.SetLevel(logrus.InfoLevel)
	Log.SetFormatter(&compactFormatter{})
	Log.SetOutput(os.Stdout)
}

// ============================================================================
// Initialization functions
// ============================================================================

// Init initializes the global logger
// If config is nil, uses default configuration (console output, info level)
func Init(cfg *Config) error {
	Log = logrus.New()

	// Use default values if no config provided
	if cfg == nil {
		cfg = &Config{Level: "info"}
	}

	// Set default values
	cfg.SetDefaults()

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("配置验证失败: %w", err)
	}

	// Set log level
	level, err := logrus.ParseLevel(cfg.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	Log.SetLevel(level)

	// Set compact formatter
	Log.SetFormatter(&compactFormatter{})
	Log.SetOutput(os.Stdout)
	Log.SetReportCaller(true)

	// Initialize Telegram Hook if configured
	if cfg.Telegram != nil && cfg.Telegram.Enabled {
		hook, err := NewTelegramHook(cfg.Telegram)
		if err != nil {
			Log.Warnf("⚠️  Failed to initialize Telegram hook: %v", err)
		} else {
			Log.AddHook(hook)
			telegramHook = hook
			minLevel := cfg.Telegram.MinLevel
			if minLevel == "" {
				minLevel = "error"
			}
			Log.Infof("✅ Telegram notification enabled (min_level: %s, chat_id: %d)", minLevel, cfg.Telegram.ChatID)
		}
	} else if cfg.Telegram != nil && !cfg.Telegram.Enabled {
		Log.Info("ℹ️  Telegram notification disabled in configuration")
	}

	return nil
}

// InitWithSimpleConfig initializes logger with simplified config
// Suitable for scenarios that only need basic functionality
func InitWithSimpleConfig(level string) error {
	return Init(&Config{Level: level})
}

// InitFromEnv 从环境变量初始化logger（包括Telegram配置）
func InitFromEnv() error {
	cfg := &Config{
		Level: os.Getenv("LOG_LEVEL"),
	}

	// Load Telegram config from environment
	telegramEnabled := os.Getenv("TELEGRAM_ENABLED")
	
	// 调试信息：显示读取到的环境变量
	fmt.Printf("[DEBUG] TELEGRAM_ENABLED=%q\n", telegramEnabled)
	
	if telegramEnabled == "true" {
		botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
		chatIDStr := os.Getenv("TELEGRAM_CHAT_ID")
		chatID := parseChatID(chatIDStr)
		minLevel := os.Getenv("TELEGRAM_MIN_LEVEL")
		
		fmt.Printf("[DEBUG] TELEGRAM_BOT_TOKEN=%s\n", maskToken(botToken))
		fmt.Printf("[DEBUG] TELEGRAM_CHAT_ID=%q (parsed: %d)\n", chatIDStr, chatID)
		fmt.Printf("[DEBUG] TELEGRAM_MIN_LEVEL=%q\n", minLevel)

		// 只有当必要字段不为空时才创建TelegramConfig
		if botToken != "" && chatID != 0 {
			telegramCfg := &TelegramConfig{
				Enabled:  true,
				BotToken: botToken,
				ChatID:   chatID,
				MinLevel: minLevel,
			}
			// 设置默认值（确保MinLevel有默认值）
			telegramCfg.SetDefaults()
			cfg.Telegram = telegramCfg
			fmt.Printf("[DEBUG] Telegram config created successfully\n")
		} else {
			// 配置不完整时给出详细警告
			fmt.Printf("[WARN] Telegram enabled but configuration incomplete:\n")
			fmt.Printf("       - bot_token present: %v\n", botToken != "")
			fmt.Printf("       - chat_id present: %v (value: %d)\n", chatID != 0, chatID)
			fmt.Printf("       请在 .env 中正确配置 TELEGRAM_BOT_TOKEN 和 TELEGRAM_CHAT_ID\n")
		}
	} else {
		fmt.Printf("[DEBUG] Telegram disabled (TELEGRAM_ENABLED != 'true')\n")
	}

	return Init(cfg)
}

// maskToken 隐藏 Token 的敏感部分
func maskToken(token string) string {
	if token == "" {
		return "(未设置)"
	}
	if len(token) < 10 {
		return "(已设置，但太短)"
	}
	return fmt.Sprintf("%s...%s", token[:8], token[len(token)-4:])
}

// parseChatID 解析ChatID字符串为int64
func parseChatID(s string) int64 {
	if s == "" {
		return 0
	}
	var chatID int64
	fmt.Sscanf(s, "%d", &chatID)
	return chatID
}

// Shutdown gracefully shuts down the logger
func Shutdown() {
	if telegramHook != nil {
		Info("📤 Shutting down Telegram notification...")
		telegramHook.Stop()
		Info("✅ Telegram notification stopped")
	}
}

// IsTelegramEnabled 检查Telegram通知是否已启用
func IsTelegramEnabled() bool {
	return telegramHook != nil && telegramHook.enabled
}

// GetTelegramStatus 获取Telegram通知状态（用于调试）
func GetTelegramStatus() string {
	if telegramHook == nil {
		return "未配置"
	}
	if !telegramHook.enabled {
		return "已禁用"
	}
	return fmt.Sprintf("已启用 (levels: %v)", len(telegramHook.levels))
}

// ============================================================================
// Logging functions
// ============================================================================

// WithFields creates logger entry with fields
func WithFields(fields logrus.Fields) *logrus.Entry {
	return Log.WithFields(fields)
}

// WithField creates logger entry with a single field
func WithField(key string, value interface{}) *logrus.Entry {
	return Log.WithField(key, value)
}

// add debug, info, warn
func Debug(args ...interface{}) {
	Log.Debug(args...)
}

func Info(args ...interface{}) {
	Log.Info(args...)
}

func Warn(args ...interface{}) {
	Log.Warn(args...)
}

func Debugf(format string, args ...interface{}) {
	Log.Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	Log.Infof(format, args...)
}

func Warnf(format string, args ...interface{}) {
	Log.Warnf(format, args...)
}

func Error(args ...interface{}) {
	Log.Error(args...)
}

func Errorf(format string, args ...interface{}) {
	Log.Errorf(format, args...)
}

func Fatal(args ...interface{}) {
	Log.Fatal(args...)
}

func Fatalf(format string, args ...interface{}) {
	Log.Fatalf(format, args...)
}

func Panic(args ...interface{}) {
	Log.Panic(args...)
}

func Panicf(format string, args ...interface{}) {
	Log.Panicf(format, args...)
}

// ============================================================================
// MCP Logger adapter
// ============================================================================

// MCPLogger adapter that allows MCP package to use the global logger
// Implements mcp.Logger interface
type MCPLogger struct{}

// NewMCPLogger creates MCP log adapter
func NewMCPLogger() *MCPLogger {
	return &MCPLogger{}
}

func (l *MCPLogger) Debugf(format string, args ...any) {
	Log.Debugf(format, args...)
}

func (l *MCPLogger) Infof(format string, args ...any) {
	Log.Infof(format, args...)
}

func (l *MCPLogger) Warnf(format string, args ...any) {
	Log.Warnf(format, args...)
}

func (l *MCPLogger) Errorf(format string, args ...any) {
	Log.Errorf(format, args...)
}
