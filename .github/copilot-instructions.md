# NOFX AI 交易系统开发指南

## 🎯 项目概览

NOFX 是一个**AI驱动的通用交易操作系统**，将多智能体决策、统一风控、低延迟执行和实时回测整合在一起。当前专注于加密货币市场，正在扩展到股票、期货、期权和外汇。

**技术栈**: Go (后端) + React + TypeScript (前端)

## 🏗️ 核心架构组件

系统采用**管理器模式**和**接口抽象**实现多交易所统一：

### 1. 交易生命周期流程

```
TraderManager → AutoTrader → StrategyEngine → DecisionEngine → MCP Client (AI)
      ↓             ↓              ↓                ↓               ↓
   管理所有     单个交易员    策略配置处理    上下文构建    AI决策
   trader实例   执行循环      指标/风控       调用AI       (DeepSeek/Qwen)
      ↓
Trader接口实现 (BinanceFutures/Bybit/OKX/Hyperliquid/Lighter)
      ↓
实际交易所API
```

### 2. 数据层架构

- **Store**: SQLite数据库抽象层 (`store/`)，每个实体有独立store（`trader.go`、`ai_model.go`、`backtest.go`等）
- **加密**: `crypto/crypto.go`提供AES-256加密，API密钥在数据库中加密存储，环境变量`DATA_ENCRYPTION_KEY`必须设置
- **Market**: 实时市场数据通过WebSocket (`market/websocket_client.go`)和REST API获取多时间框架K线数据

### 3. 关键目录映射

```
main.go              # 入口：初始化所有管理器，加载traders，启动API服务器
manager/             # TraderManager：管理多个trader实例生命周期
trader/              # AutoTrader核心逻辑 + Trader接口 + 各交易所实现
  ├── auto_trader.go        # 主交易循环（扫描→决策→执行）
  ├── interface.go          # Trader接口定义
  ├── binance_futures.go    # Binance实现
  ├── bybit_trader.go       # Bybit实现
  └── [交易所]_trader.go
decision/            # AI决策引擎：构建上下文、调用AI、解析响应
  ├── engine.go             # 核心决策逻辑
  ├── strategy_engine.go    # 策略配置管理
  └── prompt_manager.go     # Prompt模板加载
mcp/                 # Model Context Protocol客户端（AI集成）
  ├── client.go             # 基础HTTP客户端
  ├── deepseek_client.go    # DeepSeek特化客户端
  └── qwen_client.go        # Qwen特化客户端
backtest/            # 回测系统：模拟账户、历史数据回放
pool/                # 币池管理：AI500候选币、OI Top数据
api/                 # REST API服务器（Gin框架）
store/               # 数据持久化层
web/                 # React前端（Vite + TypeScript）
```

## 🔧 关键开发工作流

### 构建与测试

```powershell
# 后端测试
go test -v ./...
make test-backend

# 前端测试
cd web; npm run test

# 构建
go build -o nofx.exe main.go
cd web; npm run build

# 运行（需要.env配置）
./nofx.exe
```

### Docker部署

```bash
# docker-compose.yml定义两个服务：
# - nofx: 后端API (端口8080)
# - nofx-frontend: Nginx前端代理
docker-compose up -d
```

### 环境变量（必须）

在`.env`文件中配置：
```bash
DATA_ENCRYPTION_KEY=<32字节Base64密钥>  # 使用scripts/setup_encryption.sh生成
JWT_SECRET=<随机字符串>
RSA_PRIVATE_KEY=<RSA私钥PEM格式>       # 用于客户端加密
```

## 📝 代码约定与模式

### 1. 日志记录

统一使用`logger`包（基于zerolog），带Emoji前缀：
```go
logger.Info("✅ Configuration loaded")
logger.Infof("📊 Successfully fetched data for %d coins", len(coins))
logger.Warnf("⚠️  Failed to load: %v", err)
logger.Fatalf("❌ Fatal error: %v", err)
```

### 2. Trader接口实现

所有交易所必须实现`trader/interface.go`中的`Trader`接口：
```go
type Trader interface {
    GetBalance() (map[string]interface{}, error)
    GetPositions() ([]map[string]interface{}, error)
    OpenLong(symbol string, quantity float64, leverage int) (map[string]interface{}, error)
    OpenShort(symbol string, quantity float64, leverage int) (map[string]interface{}, error)
    CloseLong/CloseShort/CloseAll/CancelAllOrders...
}
```

### 3. 配置管理

- **全局配置**: `config/config.go` - 仅服务级配置（端口、JWT密钥）
- **Trader配置**: 存储在数据库中，通过`store.Trader()`访问
- **策略配置**: `store.StrategyConfig` - 包含指标、时间框架、风控参数、Prompt

### 4. AI决策流程

在`decision/engine.go`的`GetDecisions()`中：
1. 构建`Context`对象（账户、持仓、市场数据、指标、OI Top数据）
2. 加载Prompt模板（`prompts/`目录）
3. 调用`mcpClient.Chat()`获取AI响应
4. 使用正则表达式解析JSON（支持`<reasoning>`和`<decision>`标签）
5. 返回`[]decision.Decision`数组

**关键点**: AI响应必须是JSON数组格式，支持多种提取模式（见`engine.go`中的`reJSONFence`、`reDecisionTag`等正则）

### 5. 错误处理

- 使用标准Go error返回
- 关键操作失败记录日志但不panic（交易循环应该持续）
- 网络错误自动重试（MCP客户端有内置重试逻辑）

## 🧪 测试策略

- **单元测试**: 使用`testify`和`gomonkey`（mock），命名为`*_test.go`
- **前端测试**: Vitest框架，位于`web/src/**/*.test.tsx`
- **集成测试**: 标记为`*_integration_test.go`（如`prompt_reload_integration_test.go`）
- **竞态检测**: 关键并发代码有`*_race_test.go`

## 🔐 安全注意事项

1. **API密钥加密**: 所有敏感数据通过`crypto.CryptoService`加密后存储
2. **客户端加密**: 前端通过RSA公钥加密敏感输入，后端解密（`api/crypto_handler.go`）
3. **认证**: JWT token认证（`auth/auth.go`），中间件在`api/server.go`中
4. **2FA**: 支持OTP两步验证（`pquerna/otp`库）

## 🎨 前端开发要点

- **状态管理**: Zustand (`web/src/store/`)
- **API通信**: Axios + SWR缓存（`web/src/lib/api.ts`）
- **路由**: React Router v7（`web/src/App.tsx`）
- **样式**: Tailwind CSS + shadcn/ui组件
- **图表**: Recharts库（`web/src/components/ui/chart.tsx`）

## 🚀 添加新交易所步骤

1. 在`trader/`创建`[exchange]_trader.go`
2. 实现`Trader`接口的所有方法
3. 在`trader/auto_trader.go`的`NewTrader()`中添加switch case
4. 在`api/server.go`的`handleGetSupportedExchanges()`添加交易所信息
5. 更新前端`web/src/lib/config.ts`的交易所列表
6. 添加对应的`[exchange]_trader_test.go`

## 📚 重要文档

- **Prompt编写**: `docs/prompt-guide.zh-CN.md` - AI策略编写完整指南
- **加密系统**: `ENCRYPTION_README.md` - 端到端加密部署
- **Git工作流**: `docs/Git工作流规范.md` - 分支管理和PR流程
- **API文档**: `docs/api/` - REST API规范

## ⚡ 性能优化模式

- **缓存**: TraderManager使用`CompetitionCache`减少重复计算
- **并发**: 使用goroutine并发获取多个trader数据（加读锁）
- **池化**: `pool/coin_pool.go`缓存AI500币池数据30分钟
- **WebSocket**: 实时价格通过WS而非轮询REST

## 🐛 常见陷阱

1. **Trader并发**: 访问`TraderManager.traders`必须加锁（`RWMutex`）
2. **AI响应解析**: 支持多种格式（纯JSON、```json代码块、XML标签），解析逻辑在`decision/engine.go`
3. **时间框架**: 策略支持多时间框架（5m/15m/1h/4h），在`StrategyConfig.Timeframes`中配置
4. **Position Sync**: `trader/position_sync.go`自动同步交易所实际持仓到数据库
5. **数据库迁移**: SQLite schema变更需要添加到`migrations/`目录

## 🔄 Git分支策略

- **main**: 稳定生产分支
- **dev**: 高频开发分支（从dev创建feat/hotfix分支）
- **特性分支**: `feat/功能描述`
- **修复分支**: `hotfix/问题描述`

提交前运行`make test`确保测试通过。
