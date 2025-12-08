package logger

import (
	"fmt"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// TelegramHook 实现logrus.Hook接口，将日志推送到Telegram
type TelegramHook struct {
	sender          *TelegramSender
	levels          []logrus.Level
	enabled         bool
	lastMessageHash string            // 上一条消息的哈希（用于去重）
	lastMessageTime time.Time         // 上一条消息的时间
	duplicateCount  int               // 重复消息计数
	mu              sync.Mutex        // 保护并发访问
}

// NewTelegramHook 创建Telegram Hook
func NewTelegramHook(config *TelegramConfig) (*TelegramHook, error) {
	// 验证config不为nil
	if config == nil {
		return nil, fmt.Errorf("telegram配置不能为nil")
	}

	if !config.Enabled {
		return &TelegramHook{enabled: false}, nil
	}

	// 验证配置
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("telegram配置验证失败: %w", err)
	}

	// 确保inLevel有默认值
	config.SetDefaults()

	// 创建发送器（使用默认参数）
	sender, err := NewTelegramSender(config.BotToken, config.ChatID)
	if err != nil {
		return nil, fmt.Errorf("创建telegram发送器失败: %w", err)
	}

	hook := &TelegramHook{
		sender:          sender,
		levels:          config.GetLogrusLevels(),
		enabled:         true,
		lastMessageTime: time.Now(),
	}

	return hook, nil
}

// Levels 返回需要触发的日志级别
func (h *TelegramHook) Levels() []logrus.Level {
	if !h.enabled {
		return []logrus.Level{}
	}
	return h.levels
}

// Fire 当日志触发时调用
func (h *TelegramHook) Fire(entry *logrus.Entry) error {
	if !h.enabled {
		return nil
	}

	// 格式化消息
	message := h.formatMessage(entry)

	// 消息去重（避免短时间内发送相同消息）
	h.mu.Lock()
	messageHash := h.hashMessage(message)
	now := time.Now()

	// 如果是相同消息且30秒内，只增加计数
	if messageHash == h.lastMessageHash && now.Sub(h.lastMessageTime) < 30*time.Second {
		h.duplicateCount++
		h.mu.Unlock()
		return nil
	}

	// 如果之前有重复消息，发送汇总
	if h.duplicateCount > 0 {
		summary := fmt.Sprintf("♻️ 上一条消息在 %d 秒内重复了 %d 次", 
			int(now.Sub(h.lastMessageTime).Seconds()), h.duplicateCount)
		h.sender.SendAsync(summary)
		h.duplicateCount = 0
	}

	// 更新状态
	h.lastMessageHash = messageHash
	h.lastMessageTime = now
	h.mu.Unlock()

	// 异步发送（非阻塞）
	h.sender.SendAsync(message)

	return nil
}

// hashMessage 生成消息的简单哈希值
func (h *TelegramHook) hashMessage(message string) string {
	// 使用前100个字符作为哈希（简单且快速）
	if len(message) > 100 {
		return message[:100]
	}
	return message
}

// formatMessage 格式化日志消息为Telegram格式
func (h *TelegramHook) formatMessage(entry *logrus.Entry) string {
	// 级别emoji
	levelEmoji := h.getLevelEmoji(entry.Level)

	// 基本信息
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("%s *NOFX %s*\n", levelEmoji, strings.ToUpper(entry.Level.String())))
	builder.WriteString("\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\n")
	builder.WriteString(fmt.Sprintf("📝 *消息*: %s\n\n", escapeMarkdown(entry.Message)))

	// 字段信息（按类别分组显示）
	if len(entry.Data) > 0 {
		// Trader相关信息
		if traderID, ok := entry.Data["trader_id"]; ok {
			builder.WriteString("🤖 *Trader 信息*:\n")
			builder.WriteString(fmt.Sprintf("  • ID: `%v`\n", traderID))
			if name, ok := entry.Data["trader_name"]; ok {
				builder.WriteString(fmt.Sprintf("  • 名称: `%v`\n", name))
			}
			if exchange, ok := entry.Data["exchange"]; ok {
				builder.WriteString(fmt.Sprintf("  • 交易所: `%v`\n", exchange))
			}
			if symbol, ok := entry.Data["symbol"]; ok {
				builder.WriteString(fmt.Sprintf("  • 交易对: `%v`\n", symbol))
			}
			builder.WriteString("\n")
		}

		// 错误信息
		if err, ok := entry.Data["error"]; ok {
			builder.WriteString("⚠️ *错误详情*:\n")
			builder.WriteString(fmt.Sprintf("  `%v`\n\n", escapeMarkdown(fmt.Sprintf("%v", err))))
		}

		// 其他字段
		hasOtherFields := false
		for key, value := range entry.Data {
			if key != "trader_id" && key != "trader_name" && key != "exchange" && key != "symbol" && key != "error" {
				if !hasOtherFields {
					builder.WriteString("📊 *其他信息*:\n")
					hasOtherFields = true
				}
				builder.WriteString(fmt.Sprintf("  • %s: `%v`\n", key, value))
			}
		}
		if hasOtherFields {
			builder.WriteString("\n")
		}
	}

	// 调用位置
	if entry.HasCaller() {
		file := entry.Caller.File
		// 只保留相对路径
		if idx := strings.Index(file, "nofx/"); idx >= 0 {
			file = file[idx:]
		}
		builder.WriteString(fmt.Sprintf("📍 `%s:%d`\n", file, entry.Caller.Line))
	} else {
		// 如果entry没有caller，手动获取
		if _, file, line, ok := runtime.Caller(8); ok {
			if idx := strings.Index(file, "nofx/"); idx >= 0 {
				file = file[idx:]
			}
			builder.WriteString(fmt.Sprintf("📍 `%s:%d`\n", file, line))
		}
	}

	// 时间戳
	builder.WriteString(fmt.Sprintf("🕐 时间: `%s`", entry.Time.Format("2006-01-02 15:04:05")))

	return builder.String()
}

// getLevelEmoji 获取日志级别对应的emoji
func (h *TelegramHook) getLevelEmoji(level logrus.Level) string {
	switch level {
	case logrus.PanicLevel:
		return "🔴"
	case logrus.FatalLevel:
		return "🔴"
	case logrus.ErrorLevel:
		return "🟠"
	case logrus.WarnLevel:
		return "🟡"
	case logrus.InfoLevel:
		return "🟢"
	case logrus.DebugLevel:
		return "🔵"
	default:
		return "⚪"
	}
}

// escapeMarkdown 转义Markdown特殊字符
func escapeMarkdown(text string) string {
	replacer := strings.NewReplacer(
		"_", "\\_",
		"*", "\\*",
		"[", "\\[",
		"]", "\\]",
		"(", "\\(",
		")", "\\)",
		"~", "\\~",
		"`", "\\`",
		">", "\\>",
		"#", "\\#",
		"+", "\\+",
		"-", "\\-",
		"=", "\\=",
		"|", "\\|",
		"{", "\\{",
		"}", "\\}",
		".", "\\.",
		"!", "\\!",
	)
	return replacer.Replace(text)
}

// Stop 停止Hook（优雅关闭）
func (h *TelegramHook) Stop() {
	if h.enabled && h.sender != nil {
		// 发送一条关闭通知（可选）
		h.sender.SendAsync("📴 *NOFX System*\n━━━━━━━━━━━━━━━━\n系统正在关闭，Telegram通知服务已停止")
		// 等待消息发送完成后再关闭
		h.sender.Stop()
	}
}
