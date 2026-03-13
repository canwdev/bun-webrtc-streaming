/**
 * WebRTC 信令服务器 (Bun)
 * 提供静态页面 + WebSocket 信令，用于交换 SDP/ICE
 */
const STATIC_DIR = "./public";

const clients = new Map(); // id -> { role: 'broadcaster'|'viewer', ws }

function broadcast(fromId, message) {
  const from = clients.get(fromId);
  if (!from) return;
  const payload = JSON.stringify({ from: fromId, ...message });
  for (const [id, client] of clients) {
    if (id !== fromId && client.ws.readyState === 1) client.ws.send(payload);
  }
}

const server = Bun.serve({
  port: process.env.PORT || 3661,
  async fetch(req, server) {
    if (server.upgrade(req)) return undefined;
    const url = new URL(req.url);
    let filePath = "";
    if (url.pathname === "/" || url.pathname === "/index.html") filePath = STATIC_DIR + "/index.html";
    else if (url.pathname === "/broadcaster") filePath = STATIC_DIR + "/broadcaster.html";
    else if (url.pathname === "/viewer") filePath = STATIC_DIR + "/viewer.html";
    else if (url.pathname.startsWith("/")) filePath = (STATIC_DIR + url.pathname).replace(/\/+/g, "/");
    if (filePath) {
      const file = Bun.file(filePath);
      if (await file.exists()) return new Response(file);
    }
    return new Response("Not Found", { status: 404 });
  },
  websocket: {
    open(ws) {
      const id = crypto.randomUUID();
      ws.id = id;
      clients.set(id, { ws });
    },
    message(ws, raw) {
      try {
        const msg = JSON.parse(raw.toString());
        const id = ws.id;
        if (msg.role) {
          const c = clients.get(id);
          if (c) c.role = msg.role;
        }
        if (msg.to) {
          const target = clients.get(msg.to);
          if (target?.ws.readyState === 1) target.ws.send(JSON.stringify({ from: id, ...msg }));
        } else {
          broadcast(id, msg);
        }
      } catch (_) {}
    },
    close(ws) {
      clients.delete(ws.id);
    },
  },
});

console.log(`Server: http://localhost:${server.port}`);
console.log("  Broadcaster: http://localhost:%s/broadcaster", server.port);
console.log("  Viewer:      http://localhost:%s/viewer", server.port);
