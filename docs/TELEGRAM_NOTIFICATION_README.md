# Telegram 日志通知功能

## ✅ 已完成功能

### 核心功能
- ✅ 实时日志推送到Telegram
- ✅ 可配置日志级别（debug/info/warn/error/fatal/panic）
- ✅ 异步非阻塞发送
- ✅ 自动重试机制（失败后重试3次）
- ✅ 消息去重（30秒内相同消息自动合并）
- ✅ 优雅关闭支持
- ✅ 美化的Markdown格式消息

### 文件清单

**核心实现**:
- `logger/telegram_hook.go` - Logrus Hook实现
- `logger/telegram_sender.go` - 异步消息发送器
- `logger/config.go` - 配置结构
- `logger/logger.go` - 集成到主Logger

**文档**:
- `docs/TELEGRAM_NOTIFICATION_GUIDE.md` - 完整使用指南
- `.env.example` - 环境变量配置示例

**测试**:
- `test_telegram.go` - 功能测试脚本

## 🚀 快速使用

### 1. 配置环境变量（`.env`）

```bash
TELEGRAM_ENABLED=true
TELEGRAM_BOT_TOKEN=your-bot-token-from-BotFather
TELEGRAM_CHAT_ID=your-chat-id
TELEGRAM_MIN_LEVEL=error
```

### 2. 运行测试

```bash
# 测试Telegram通知功能
go run test_telegram.go

# 正常启动系统
go run main.go
```

### 3. 验证

检查控制台输出是否有：
```
✅ Telegram notification enabled (min_level: error)
```

然后触发错误日志，查看Telegram是否收到通知。

## 📋 配置说明

### 日志级别

| 级别 | 说明 | 使用场景 |
|------|------|---------|
| `error` | 错误及以上 | **推荐**（生产环境） |
| `warn` | 警告及以上 | 测试环境 |
| `info` | 信息及以上 | 需详细监控 |
| `debug` | 所有日志 | 仅调试使用 |

### 消息格式

```
🟠 NOFX ERROR
────────────────
📝 消息: API request failed

🤖 Trader 信息:
  • ID: trader-001
  • 名称: BTC-Strategy
  • 交易所: binance
  • 交易对: BTCUSDT

⚠️ 错误详情:
  connection timeout

📍 trader/auto_trader.go:123
🕐 时间: 2025-12-08 15:30:45
```

## 🔧 技术细节

### 架构

```
Logger 
  ↓ (AddHook)
TelegramHook
  ↓ (Fire)
formatMessage() → 格式化日志
  ↓
SendAsync() → 写入缓冲区（非阻塞）
  ↓
listenAndSend() → goroutine监听
  ↓
sendWithRetry() → 发送到Telegram API（带重试）
```

### 性能特性

- **非阻塞**: 使用channel异步发送，不影响主程序
- **缓冲区**: 20条消息缓冲，超出则丢弃（避免内存溢出）
- **去重**: 30秒内相同消息合并为汇总
- **重试**: 失败后每3秒重试，最多3次

### 依赖

```go
github.com/go-telegram-bot-api/telegram-bot-api/v5
github.com/sirupsen/logrus
```

## 📚 相关文档

- [完整使用指南](./TELEGRAM_NOTIFICATION_GUIDE.md) - 包含配置、测试、故障排查
- [Telegram Bot API](https://core.telegram.org/bots/api) - 官方API文档

## 🐛 常见问题

**Q: 没有收到消息？**
- 检查Bot Token和Chat ID是否正确
- 确认已给Bot发送 `/start` 命令
- 检查日志级别配置（如设置error但只触发info日志）

**Q: 消息太多怎么办？**
- 调高 `TELEGRAM_MIN_LEVEL`（如设置为 `error` 或 `fatal`）
- 系统已内置消息去重

**Q: 如何测试配置？**
```bash
go run test_telegram.go
```

## 🎯 使用建议

1. **生产环境**: `TELEGRAM_MIN_LEVEL=error`
2. **测试环境**: `TELEGRAM_MIN_LEVEL=warn`
3. **开发环境**: `TELEGRAM_MIN_LEVEL=info` 或不启用
4. **团队协作**: 使用群组Chat ID
5. **安全**: 不要将Bot Token提交到Git

---

更新时间: 2025-12-08
