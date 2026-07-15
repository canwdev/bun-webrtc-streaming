# WebRTC 摄像头直播

使用 WebRTC + 原生 HTML5 实现的 WebCam 低延迟直播，支持 **Bun** 和 **Go** 两种后端。

## 功能

- **极低延迟**：目标 100ms 以内（本地/局域网下可达）
- **网页调节**：分辨率（320×240～1920×1080）、帧率（15/24/30/60 fps）
- **摄像头控制**：选择设备、开始/停止推流
- **双后端支持**：Bun (Node.js 生态) 或 Go (单文件 exe)

## 项目结构

```
├── public/                # 共享前端静态文件
│   ├── index.html         # 首页
│   ├── broadcaster.html   # 推流端
│   └── viewer.html        # 观看端
├── nodejs/                # Bun 后端
│   ├── server.js          # 信令服务器
│   ├── package.json
│   └── README.md
├── go/                    # Go 后端
│   ├── main.go            # 信令服务器（HTML 内嵌）
│   ├── go.mod / go.sum
│   ├── package.json       # Bun 构建脚本（dev/build/cross-compile）
│   ├── .air.toml          # 热重载配置
│   ├── build.bat          # Windows 构建脚本
│   ├── bin/               # 构建产物
│   ├── public/            # 构建时从 ../public 自动复制（不纳入版本控制）
│   └── README.md
├── .github/workflows/     # CI/CD
│   └── release-go.yml     # master 推送自动构建预览版
├── package.json           # 根目录 Bun 脚本入口
└── README.md
```

## 快速开始

### 方式一：Bun 后端

```bash
cd nodejs
bun run server.js
```

详见 [nodejs/README.md](nodejs/README.md)

### 方式二：Go 后端（单文件 exe）

```bash
cd go
go build -trimpath -ldflags="-s -w" -o webrtc-server .
.\webrtc-server.exe
```

编译产物为单文件可执行程序，HTML 已内嵌，无需额外文件。

详见 [go/README.md](go/README.md)

## 访问

启动后端后打开浏览器：

- 首页：<http://localhost:3661>
- 推流端：<http://localhost:3661/broadcaster>
- 观看端：<http://localhost:3661/viewer>

## 项目工作原理

- **架构**：媒体走 WebRTC P2P（浏览器直连，不经服务器），信令走 WebSocket
- **信令服务器**：提供静态页 + WebSocket；带 `to` 则点对点转发，否则广播
- **推流端**：getUserMedia 采集；每个观看端一条 RTCPeerConnection，发 offer、收 answer/ICE
- **观看端**：收到 offer 后建 PC、回 answer、交换 ICE，播放收到的流

## 低延迟说明

- 端到端为 WebRTC P2P，无服务器转码，延迟主要来自采集、编码、网络与解码
- 本机或局域网下更容易达到 100ms 内
- 若需公网低延迟，需部署 TURN/STUN 或媒体服务器

## License

MIT
