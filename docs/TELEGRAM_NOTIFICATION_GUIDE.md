# Telegram 日志通知功能使用指南

## 📱 功能概述

Telegram日志通知功能可以将NOFX系统的关键日志实时推送到你的Telegram，让你随时随地监控系统运行状态。

### ✨ 主要特性

- 🔔 **实时推送**: 关键错误、警告等日志立即通知
- 🎚️ **可配置级别**: 灵活设置推送的最低日志级别
- 🚀 **异步非阻塞**: 不影响系统性能
- 🔄 **自动重试**: 网络故障时自动重试发送
- 📊 **格式化展示**: 清晰的Markdown格式，包含Trader信息、错误详情、调用位置等
- ♻️ **消息去重**: 30秒内相同消息自动合并，避免刷屏

## 🚀 快速开始

### 1️⃣ 创建Telegram Bot

1. 在Telegram搜索 `@BotFather`
2. 发送 `/newbot` 命令
3. 按提示设置Bot名称和用户名
4. 保存Bot Token（格式：`123456789:ABCdefGHIjklMNOpqrsTUVwxyz`）

### 2️⃣ 获取Chat ID

**方法一：使用 @userinfobot**
1. 在Telegram搜索 `@userinfobot`
2. 发送任意消息
3. 获取你的Chat ID（纯数字）

**方法二：使用 @getidsbot**
1. 搜索并启动 `@getidsbot`
2. 发送 `/start`
3. 复制显示的ID

**方法三：通过API获取（群组）**
1. 将Bot添加到群组
2. 在群组发送消息
3. 访问：`https://api.telegram.org/bot<YOUR_BOT_TOKEN>/getUpdates`
4. 在返回的JSON中查找 `"chat":{"id":-1001234567890}` 

### 3️⃣ 配置环境变量

在 `.env` 文件中添加：

```bash
# 开启Telegram通知
TELEGRAM_ENABLED=true

# Bot Token（从 @BotFather 获取）
TELEGRAM_BOT_TOKEN=123456789:ABCdefGHIjklMNOpqrsTUVwxyz

# Chat ID（个人或群组ID）
TELEGRAM_CHAT_ID=123456789

# 最低日志级别（可选，默认error）
# 选项: debug, info, warn, error, fatal, panic
TELEGRAM_MIN_LEVEL=error
```

### 4️⃣ 重启系统

```bash
# 停止正在运行的实例
# 重新启动
go run main.go

# 或 Docker
docker-compose restart nofx
```

## 📋 配置说明

### 日志级别说明

| 级别 | 说明 | 建议场景 |
|------|------|---------|
| `panic` | 仅系统崩溃 | 不推荐（太严格） |
| `fatal` | 致命错误 | 生产环境最严格模式 |
| `error` | 错误信息 | **推荐**（默认） |
| `warn` | 警告信息 | 开发/测试环境 |
| `info` | 一般信息 | 需要详细监控时 |
| `debug` | 调试信息 | 仅调试时使用（消息量大） |

**推荐配置**:
- **生产环境**: `error` 或 `fatal`
- **测试环境**: `warn`
- **开发环境**: `info` 或 `debug`

### 消息格式示例

#### 错误消息

```
🟠 NOFX ERROR
────────────────────────────
📝 消息: Failed to execute trade order

🤖 Trader 信息:
  • ID: `trader-001`
  • 名称: `BTC-Long-Strategy`
  • 交易所: `binance`
  • 交易对: `BTCUSDT`

⚠️ 错误详情:
  `insufficient balance: required 1000 USDT, available 500 USDT`

📍 `trader/binance_futures.go:245`
🕐 时间: `2025-12-08 15:30:45`
```

#### 警告消息

```
🟡 NOFX WARN
────────────────────────────
📝 消息: API rate limit approaching

📊 其他信息:
  • requests_remaining: `5`
  • reset_time: `60s`

📍 `market/api_client.go:123`
🕐 时间: `2025-12-08 15:30:45`
```

## 🔧 高级功能

### 1. 消息去重

系统自动检测30秒内的重复消息，避免刷屏：

```
原始消息: "Database connection failed"
（30秒内又发生3次相同错误）
汇总消息: "♻️ 上一条消息在 30 秒内重复了 3 次"
```

### 2. 异步发送

- 消息缓冲区大小：20条
- 超出缓冲区时自动丢弃（不阻塞主程序）
- 发送失败自动重试3次，间隔3秒

### 3. 群组通知

将Bot添加到Telegram群组可实现团队协作：

1. 创建群组
2. 将Bot添加为成员
3. 使用群组的Chat ID（负数，如 `-1001234567890`）

### 4. 多Bot配置

如需不同级别的日志发送到不同Bot：

```bash
# 方案：使用多个实例，每个配置不同的环境变量
# 或：在代码中手动添加多个TelegramHook
```

## 🧪 测试

### 测试Telegram配置

创建测试脚本 `test_telegram.go`:

```go
package main

import (
	"nofx/logger"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	logger.InitFromEnv()
	
	logger.Info("🟢 这是INFO级别测试")
	logger.Warn("🟡 这是WARN级别测试")
	logger.Error("🟠 这是ERROR级别测试")
	
	logger.Shutdown()
}
```

运行测试：
```bash
go run test_telegram.go
```

### 验证配置

启动系统后，查看日志输出：

```
✅ Telegram notification enabled (min_level: error)
```

如果看到此消息，说明配置成功。

## ⚠️ 注意事项

### 1. Bot Token 安全

- ⛔ **永远不要**将Bot Token提交到Git仓库
- ✅ 使用 `.env` 文件（已在 `.gitignore` 中）
- ✅ 生产环境使用环境变量或密钥管理服务

### 2. 消息频率

- Telegram限制：每秒最多30条消息/Bot
- 系统内置缓冲区和去重机制
- 建议日志级别设置为 `error` 或更高

### 3. Chat ID 类型

- 个人ID：正数（如 `123456789`）
- 群组ID：负数（如 `-1001234567890`）
- 频道ID：以 `-100` 开头

### 4. 网络问题

- 国内网络可能需要代理才能访问Telegram API
- 设置代理：`export HTTP_PROXY=http://127.0.0.1:7890`
- 或使用Telegram Bot API本地服务器

## 🐛 故障排查

### 问题1: 没有收到消息

**检查清单**:
- [ ] Bot Token是否正确
- [ ] Chat ID是否正确（个人/群组）
- [ ] 是否已给Bot发送过 `/start` 命令
- [ ] 日志级别配置是否合理（如设置为error但触发的是info日志）
- [ ] 网络是否可以访问Telegram API

**调试步骤**:
```bash
# 1. 测试Bot Token
curl "https://api.telegram.org/bot<YOUR_TOKEN>/getMe"

# 2. 查看系统日志
tail -f logs/nofx.log | grep Telegram

# 3. 手动发送测试消息
curl -X POST "https://api.telegram.org/bot<YOUR_TOKEN>/sendMessage" \
  -d "chat_id=<YOUR_CHAT_ID>" \
  -d "text=Test message"
```

### 问题2: 消息格式错乱

**原因**: Markdown转义问题

**解决方案**: 
- 系统已内置转义功能
- 如仍有问题，检查日志消息中是否包含特殊字符

### 问题3: 发送失败错误

常见错误码：

| 错误 | 原因 | 解决方案 |
|------|------|---------|
| 400 Bad Request | Chat ID错误 | 验证Chat ID格式 |
| 401 Unauthorized | Token错误 | 检查Bot Token |
| 403 Forbidden | Bot被阻止 | 重新启动Bot（发送 `/start`） |
| 429 Too Many Requests | 超过频率限制 | 降低日志级别 |

### 问题4: 系统启动慢

如果Telegram配置错误导致启动变慢：
- 设置 `TELEGRAM_ENABLED=false` 暂时禁用
- 或修复配置后重启

## 📊 性能影响

- **CPU使用**: 几乎无影响（异步发送）
- **内存占用**: ~2MB（缓冲区+goroutine）
- **网络带宽**: 取决于消息频率，每条消息 ~1KB

## 🔄 更新日志

### v1.0 (2025-12-08)
- ✅ 基础Telegram推送功能
- ✅ 支持自定义日志级别
- ✅ 消息去重机制
- ✅ 异步发送和自动重试
- ✅ 优雅关闭支持
- ✅ 美化的消息格式

## 📚 相关资源

- [Telegram Bot API 文档](https://core.telegram.org/bots/api)
- [@BotFather](https://t.me/BotFather) - 官方Bot创建工具
- [go-telegram-bot-api](https://github.com/go-telegram-bot-api/telegram-bot-api) - Go SDK

## 💡 最佳实践

1. **生产环境**: 设置 `TELEGRAM_MIN_LEVEL=error`，只接收重要通知
2. **团队协作**: 使用群组，让团队成员都能收到通知
3. **分级通知**: 重要系统用 `error`，测试环境用 `warn`
4. **定期检查**: 确保Bot Token未过期，群组ID未变更
5. **日志审查**: 定期查看Telegram消息，及时处理异常

---

如有问题或建议，请提交Issue或PR！🚀
