# Go 后端

使用 Go + gorilla/websocket 的信令服务器，编译为单文件可执行程序，HTML 静态文件内嵌。

## 环境要求

- [Go](https://go.dev/dl/) >= 1.21
- [Bun](https://bun.sh/)（用于 `bun run` 构建脚本，非运行时必需）
- [air](https://github.com/air-verse/air)（可选，用于热重载开发）

```bash
go install github.com/air-verse/air@latest
```

## 快速开始

### 开发模式（热重载）

```bash
cd go
bun install
bun run dev
```

`air` 会监听 `.go` 和 `.html` 文件变更，自动重新编译并重启服务。

### 构建当前平台

```bash
cd go
bun run build
```

产物输出到 `go/bin/webrtc-server`（或 `.exe`）。

### 跨平台交叉编译（全部 6 个目标）

```bash
cd go
bun install
bun run build:all
```

产物输出到：

```
go/bin/
├── windows_amd64/webrtc-server.exe
├── windows_arm64/webrtc-server.exe
├── linux_amd64/webrtc-server
├── linux_arm64/webrtc-server
├── mac_amd64/webrtc-server
└── mac_arm64/webrtc-server
```

也可单独构建某个平台：

```bash
bun run build:win:amd64
bun run build:linux:amd64
bun run build:mac:arm64
```

### 手动构建（无需 Bun）

```bash
cd go
# 先复制前端文件
cp -r ../public public    # Linux/macOS
xcopy /Y /Q ..\public\* public\    # Windows

# 编译
go build -trimpath -ldflags="-s -w" -o webrtc-server .
```

编译产物为**单文件可执行程序**，HTML 内嵌，无需额外文件即可运行。

## 运行

```bash
.\webrtc-server.exe
```

### 自定义端口

```bash
# PowerShell
$env:PORT="8080"; .\webrtc-server.exe

# Linux/macOS
PORT=8080 ./webrtc-server
```

默认端口：`3661`

## 访问

启动后打开浏览器访问：

- 首页：<http://localhost:3661>
- 推流端：<http://localhost:3661/broadcaster>
- 观看端：<http://localhost:3661/viewer>

## 构建脚本说明

| 命令 | 说明 |
|------|------|
| `bun run dev` | 开发模式，使用 air 热重载 |
| `bun run build` | 构建当前平台 |
| `bun run build:all` | 交叉编译全部 6 个平台 |
| `bun run build:win:amd64` | Windows x64 |
| `bun run build:win:arm64` | Windows ARM64 |
| `bun run build:linux:amd64` | Linux x64 |
| `bun run build:linux:arm64` | Linux ARM64 |
| `bun run build:mac:amd64` | macOS Intel |
| `bun run build:mac:arm64` | macOS Apple Silicon |
| `bun run clean` | 清理 bin/ 目录 |
| `bun run copy:frontend` | 从 ../public 复制前端文件 |

## 目录结构

```
go/
├── main.go           # 信令服务器（HTTP + WebSocket，HTML 内嵌）
├── go.mod / go.sum   # Go 依赖
├── package.json      # Bun 构建脚本
├── .air.toml         # air 热重载配置
├── .gitignore
├── build.bat         # Windows 构建脚本（纯批处理，无需 Bun）
├── README.md
├── bin/              # 构建产物（不纳入版本控制）
├── tmp/              # 开发临时文件（不纳入版本控制）
└── public/           # 构建时从 ../public 自动复制（不纳入版本控制）
```
