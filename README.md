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

## 项目工作原理

- **目标**：用 Bun + WebRTC + 原生 HTML5 做摄像头直播，推流端采集，多观看端实时看，低延迟。
- **架构**：媒体走 WebRTC P2P（浏览器直连，不经服务器），信令走 WebSocket（经 server.js 转发 SDP/ICE）。
- **信令服务器**：提供静态页 + WebSocket；带 `to` 则点对点转发，否则广播；消息带 `from` 标识发送方。
- **推流端**：getUserMedia 采集；每收到一个 `role: "viewer"` 就为该观看端建一条 RTCPeerConnection，发 offer、收 answer/ICE；同一路流通过多条 PC 推给多个观看端。
- **观看端**：发 `role: "viewer"`；收到 `broadcaster_ready` 再发一次以支持「先开观看端」；收到 offer 后建 PC、回 answer、交换 ICE，把收到的流播在 `<video>` 上。

## 接收端分辨率从低到高是谁控制的？

- **采集**：`getUserMedia` 用 `ideal` 时，设备可能先给低分辨率再提升；用 `exact` 可强制首帧即目标分辨率（部分设备可能不支持）。
- **编码**：Chrome 等会先以较低分辨率/码率编码再随带宽探测逐步拉高，无法用 JS 关闭。
- 接收端只是显示「当前收到的编码分辨率」，从低到高由推流端采集 + 浏览器编码器共同决定。

## 多个观看端是怎么推流的？

- 推流端为**每个观看端**维护**一条** RTCPeerConnection（存在 `peerConnections` Map 里）。
- 同一路 `stream`（getUserMedia）通过 `addTrack(track, stream)` 加到每一条 PC 上。
- 即：**N 个观看端 = N 条独立的 WebRTC 连接**，每条都传同一路摄像头流；媒体是推流端 ↔ 各观看端点对点，不经服务器。
- 编码/发送是「每条连接一份」，观看端越多，推流端上行和 CPU 压力越大。

## 纯局域网、不连外网能用吗？

- **可以。**
- 信令：用局域网 IP 访问服务器（如 `http://192.168.1.100:3661`）即可，不需要外网。
- WebRTC：当前未配 STUN/TURN（`iceServers` 为空），ICE 用局域网 host 候选即可在同网段直连，不依赖外网。
- 使用：在同一局域网内，推流端和观看端都通过该机器的局域网 IP 打开对应页面即可。

## License

MIT
