package main

import (
	"fmt"
	"nofx/logger"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	// 加载环境变量
	_ = godotenv.Load()

	fmt.Println("🔍 检查 Telegram 配置...")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// 检查环境变量
	enabled := os.Getenv("TELEGRAM_ENABLED")
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	chatID := os.Getenv("TELEGRAM_CHAT_ID")
	minLevel := os.Getenv("TELEGRAM_MIN_LEVEL")

	fmt.Printf("环境变量检查:\n")
	fmt.Printf("  TELEGRAM_ENABLED: %q\n", enabled)
	fmt.Printf("  TELEGRAM_BOT_TOKEN: %s\n", maskToken(botToken))
	fmt.Printf("  TELEGRAM_CHAT_ID: %q\n", chatID)
	fmt.Printf("  TELEGRAM_MIN_LEVEL: %q (默认: error)\n", minLevel)
	fmt.Println()

	// 初始化 logger
	fmt.Println("📋 初始化 logger...")
	if err := logger.InitFromEnv(); err != nil {
		fmt.Printf("❌ 初始化失败: %v\n", err)
		return
	}

	// 检查 Telegram 状态
	status := logger.GetTelegramStatus()
	fmt.Printf("📊 Telegram 状态: %s\n", status)
	fmt.Println()

	if !logger.IsTelegramEnabled() {
		fmt.Println("⚠️  Telegram 通知未启用")
		fmt.Println()
		fmt.Println("💡 启用方法：")
		fmt.Println("   1. 在 .env 文件中添加：")
		fmt.Println("      TELEGRAM_ENABLED=true")
		fmt.Println("      TELEGRAM_BOT_TOKEN=your-bot-token")
		fmt.Println("      TELEGRAM_CHAT_ID=your-chat-id")
		fmt.Println("      TELEGRAM_MIN_LEVEL=error")
		fmt.Println()
		fmt.Println("   2. 获取 Bot Token: 联系 @BotFather")
		fmt.Println("   3. 获取 Chat ID: 联系 @userinfobot")
		return
	}

	// 测试发送消息
	fmt.Println("🧪 开始测试 Telegram 通知...")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	// 测试不同级别的日志
	testMessages := []struct {
		level   string
		logFunc func(string)
	}{
		{"ERROR", logger.Error},
		{"WARN", logger.Warn},
		{"INFO", logger.Info},
	}

	for _, tm := range testMessages {
		msg := fmt.Sprintf("🧪 测试 %s 级别消息 - %s", tm.level, time.Now().Format("15:04:05"))
		fmt.Printf("  发送: %s\n", msg)
		tm.logFunc(msg)
		time.Sleep(2 * time.Second)
	}

	fmt.Println()
	fmt.Println("⏳ 等待消息发送完成...")
	time.Sleep(5 * time.Second)

	// 关闭
	logger.Shutdown()

	fmt.Println()
	fmt.Println("✅ 测试完成！")
	fmt.Println()
	fmt.Println("📱 请检查 Telegram 是否收到消息。")
	fmt.Println()
	fmt.Println("如果没有收到消息，请检查：")
	fmt.Println("  1. Bot Token 和 Chat ID 是否正确")
	fmt.Println("  2. 是否已给 Bot 发送过 /start 命令")
	fmt.Println("  3. TELEGRAM_MIN_LEVEL 设置（只会发送该级别及以上的消息）")
	fmt.Println("  4. 网络是否可以访问 Telegram API")
}

func maskToken(token string) string {
	if token == "" {
		return "(未设置)"
	}
	if len(token) < 10 {
		return "(已设置，但太短)"
	}
	return fmt.Sprintf("%s...%s (已设置)", token[:8], token[len(token)-4:])
}
