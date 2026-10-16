# AI 助手配置

管理入口：后台「博客设置 → 互动与联系 → AI 助手」。配置保存在现有 `configs` 集合，`typ` 为 `ai_assistant`。本模块提供配置管理和简短对话测试，不创建访客会话或执行博客工具。模型接入统一使用 Eino。

## 部署密钥

模型 API Key 的加密密钥直接保存在后端启动 YAML 的 `ai_assistant.encryption_key` 中，值为 32 个随机字节的标准 Base64 编码。现有四份 `server/config/fnote*.yaml` 已配置同一份随机密钥，可直接保存模型配置，无需设置环境变量。配置结构为：

```yaml
ai_assistant:
  encryption_key: "<32 字节随机密钥的 Base64 编码>"
```

如需自行生成密钥，可以使用 Node.js：

```sh
node -e "console.log(require('crypto').randomBytes(32).toString('base64'))"
```

本地默认读取 `server/config/fnote.yaml`；通过 `--config` 指定其他文件时，在对应文件中设置。Docker 镜像默认读取 `fnote.docker.yaml`。修改后重启 Go 服务；镜像内的配置发生修改后，需要重新构建并部署镜像。

该值用于 AES-256-GCM 加密，需保持稳定，不能每次启动重新生成。已有密文恢复到其他环境时也需要同一密钥；更换后需重新填写各 Provider 的 API Key。未设置时，博客和配置读取仍然可用，但不能保存新的模型密钥或启用已有密钥。

## 管理接口

- `GET /admin-api/configs/ai-assistant`
- `PUT /admin-api/configs/ai-assistant`
- `POST /admin-api/configs/ai-assistant/test-provider`

沿用现有管理员 JWT 鉴权。读取不存在的配置返回默认关闭状态；保存时自动创建配置。PUT 提交完整配置，包含 `enabled`、`assistant_name`、`system_prompt`、`enabled_tools`、`active_provider_id`、`providers`。

Provider 字段为 `id`、`name`、`api_type`、`base_url`、`model`、`api_key`，以及可选的 `max_output_tokens`、`reasoning_effort`、`thinking_mode`。当前支持 GPT 和 DeepSeek，暂不提供 Claude 选项。

| 服务 | api_type | 默认服务地址 | 参数映射 |
| --- | --- | --- | --- |
| GPT | `openai_responses` | `https://api.openai.com/v1` | `max_output_tokens`；`reasoning_effort` 映射到 `reasoning.effort` |
| DeepSeek | `openai_chat_completions` | `https://api.deepseek.com` | 输出上限映射到 `max_tokens`；`thinking_mode` 映射到 `thinking.type`；推理强度使用 `reasoning_effort` |

留空的可选参数不发送，由服务商使用默认值。输出上限的应用范围为 16–32768，模型本身可能有更严格的限制。不同 GPT 模型支持的推理强度不同，不能只根据协议保证每个值均可用，需通过对话测试验证。DeepSeek 思考模式关闭时不能同时指定开启推理的强度。当前不默认发送 temperature、top_p 等可能与推理模式不兼容的参数。

协议、消息构造及响应解析由 Eino 的 `agenticopenai` 组件处理，GPT 使用 Responses 模型，DeepSeek 使用 OpenAI 兼容 Chat Completions 模型。共享构造入口位于 `server/internal/pkg/aimodel`，供后续助手生成复用。

读取仅返回 `has_api_key`，不会返回明文或密文。保存时同一 Provider ID 的 `api_key` 留空或省略，保留原密钥；填写新值则替换；删除 Provider 同时移除其密钥。修改服务地址时，请确认该服务可以使用保留的密钥。

工具仅允许 `search_posts`、`read_post`、`get_blog_overview`、`get_author_profile`。启用助手必须选择存在的 Provider 并配置可解密的 API Key。后台页面新增服务时，即使暂不开启助手，也需要先填写有效密钥并通过连接测试。

## 对话测试

每个模型服务提供「测试连接」按钮，测试当前表单中的地址、模型、协议、推理选项和密钥，不要求先保存或启用助手。密钥留空时按 Provider ID 读取已保存的密钥；新填密钥直接用于本次请求，不写入数据库，也不要求已配置加密主密钥。

点击后会实际发送一条「请只回复 OK。」的非流式对话请求，收到有效的助手文本后才显示成功，并显示回复和耗时。无效密钥、模型或地址、限流、网络故障、超时和空回复均显示失败。测试有 30 秒总时限，全站最多并发两个测试，不自动重试或跟随重定向。未填写输出上限时，测试使用 1024 token 的预算；有填写时使用当前填写值。测试不验证博客工具和流式能力，也不自动保存配置。

后台页面要求新增或修改连接参数的模型服务先通过「测试连接」，全部通过后才能保存；修改地址、模型、接口协议、密钥或调用参数后，对应测试结果立即失效，需要重新测试。仅修改服务名称、助手名称、提示词或工具选择不影响已有连接测试。删除服务不要求测试。已保存且连接参数未变化的服务无需重新测试。

测试通过后，保存接口仅进行配置校验、密钥加密和持久化，不重复调用模型。测试状态只用于后台页面的保存限制，不写入数据库，也不增加后端测试凭证。

页面离开时取消测试。上游错误正文不会直接返回，测试请求和响应不进入正文日志。服务成功回复中若包含本次密钥，也会先进行隐藏。

参考：[OpenAI 推理参数](https://developers.openai.com/api/docs/guides/reasoning)、[DeepSeek 思考模式](https://api-docs.deepseek.com/guides/thinking_mode/)、[Eino Agentic OpenAI 组件](https://github.com/cloudwego/eino-ext/tree/main/components/model/agenticopenai)。

本模块不增加公开配置端点，也不修改前台演示。AI 展示字段和真实聊天将在后续模块接入。
