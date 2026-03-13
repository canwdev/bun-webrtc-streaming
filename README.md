# WebRTC 摄像头直播

使用 Bun + WebRTC + 原生 HTML5 实现的 WebCam 低延迟直播。

## 功能

- **极低延迟**：目标 100ms 以内（本地/局域网下可达）
- **网页调节**：分辨率（320×240～1920×1080）、帧率（15/24/30/60 fps）
- **摄像头控制**：选择设备、开始/停止推流

## 环境

- [Bun](https://bun.sh/) 作为运行时与包管理器
- 现代浏览器（Chrome / Edge / Firefox）支持 WebRTC、`getUserMedia`

## 运行

```bash
# 安装依赖（本项目无额外依赖，仅需 Bun）
bun install

# 启动服务
bun start
```

然后：

1. 打开 **推流端**：<http://localhost:3000/broadcaster>
2. 选择摄像头、分辨率、帧率，点击「开始推流」
3. 打开 **观看端**：<http://localhost:3000/viewer>（可另开标签页或另一设备同网段访问）

## 低延迟说明

- 端到端为 WebRTC P2P，无服务器转码，延迟主要来自采集、编码、网络与解码。
- 在 **本机或局域网** 下更容易达到 100ms 内；分辨率、帧率越低，编码负担越小，通常延迟更稳定。
- 若需公网低延迟，需部署 TURN/STUN 或媒体服务器，并保证带宽与网络质量。

## 项目结构

```
├── server.js           # Bun 信令服务（HTTP + WebSocket）
├── public/
│   ├── index.html      # 首页
│   ├── broadcaster.html # 推流端（摄像头、分辨率、帧率、启停）
│   └── viewer.html     # 观看端
├── package.json
└── README.md
```

## License

MIT
