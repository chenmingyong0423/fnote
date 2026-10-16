/* eslint-disable @typescript-eslint/no-require-imports -- Test the streaming client without another build tool. */
const assert = require("node:assert/strict");
const fs = require("node:fs");
const vm = require("node:vm");
const ts = require("typescript");

const source = ts.transpileModule(fs.readFileSync("src/api/ai.ts", "utf8"), {
  compilerOptions: { target: ts.ScriptTarget.ES2022, module: ts.ModuleKind.CommonJS },
}).outputText;
const mod = { exports: {} };
vm.runInThisContext(`(function(module,exports){${source}\n})`)(mod, mod.exports);
const api = mod.exports;
const encoder = new TextEncoder();
const streamResponse = bytes => new Response(new ReadableStream({
  start(controller) { const size = bytes.length > 4096 ? 4096 : 1; for (let i = 0; i < bytes.length; i += size) controller.enqueue(bytes.slice(i, i + size)); controller.close(); },
}), { headers: { "Content-Type": "text/event-stream; charset=utf-8" } });

(async () => {
  const frames = ': heartbeat\r\n\r\nevent: meta\r\ndata: {"type":"meta","message_id":"answer"}\r\n\r\n' +
    'event: delta\r\ndata: {"type":"delta",\r\ndata: "delta":"你好🌍"}\r\n\r\n' +
    'event: done\r\ndata: {"type":"done","status":"complete"}\r\n\r\n';
  const events = [];
  await api.consumeChatStream(streamResponse(encoder.encode(frames)), event => events.push(event));
  assert.deepEqual(events.map(e => e.type), ["meta", "delta", "done"]);
  assert.equal(events[1].delta, "你好🌍", "arbitrary byte boundaries must preserve Unicode");
  await assert.rejects(api.consumeChatStream(streamResponse(encoder.encode('data: {"type":"delta","delta":"partial"}\n\n')), () => {}), /中断/);
  await assert.rejects(api.consumeChatStream(new Response("{}"), () => {}), /有效的数据流/);
  await assert.rejects(api.consumeChatStream(streamResponse(encoder.encode('data: {"type":"unknown"}\n\n')), () => {}), /格式/);
  await assert.rejects(api.consumeChatStream(streamResponse(encoder.encode("data: " + "x".repeat(128 * 1024) + "\n\n")), () => {}), /过长/);
  const original = global.fetch;
  try {
    let request;
    global.fetch = async (url, init) => { request = { url, init }; return streamResponse(encoder.encode('data: {"type":"done","status":"stopped"}\n\n')); };
    await api.sendChat("conversation", "stable-request-id", "问题", new AbortController().signal);
    assert.equal(request.url, "/api/ai/conversations/conversation/messages");
    assert.equal(request.init.credentials, "same-origin");
    assert.deepEqual(JSON.parse(request.init.body), { request_id: "stable-request-id", content: "问题" });
    global.fetch = async () => new Response('{"message":"对话不存在或无权访问"}', { status: 404 });
    await assert.rejects(api.getChatConversation("private"), /无权访问/);
  } finally { global.fetch = original; }
  console.log("PASS: SSE UTF-8/CRLF/multiline framing, heartbeats, interruption, malformed/oversized events, visitor credentials and stable request IDs");
})().catch(error => { console.error(error); process.exitCode = 1; });
