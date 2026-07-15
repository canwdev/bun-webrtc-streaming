# Bun 后端

使用 Bun 运行时 + 原生 WebSocket 的信令服务器。

## 环境要求

- [Bun](https://bun.sh/) >= 1.0

## 构建 & 运行

无需编译，直接运行：

```bash
cd nodejs
bun run server.js
```

或者使用 package.json 脚本：

```bash
cd nodejs
bun start        # 生产模式
bun run dev      # 开发模式（热重载）
```

### 自定义端口

```bash
PORT=8080 bun run server.js
```

默认端口：`3661`

## 访问

启动后打开浏览器访问：

- 首页：<http://localhost:3661>
- 推流端：<http://localhost:3661/broadcaster>
- 观看端：<http://localhost:3661/viewer>

## 目录结构

```
nodejs/
├── server.js        # 信令服务器（HTTP + WebSocket）
├── package.json     # 项目配置与脚本
├── README.md
└── ../public/       # 共享前端静态文件
    ├── index.html
    ├── broadcaster.html
    └── viewer.html
```
