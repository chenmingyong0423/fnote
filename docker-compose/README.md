# docker-compose 部署说明

当前目录提供两套部署方式：

- `docker-compose.yaml`：本地直连部署，`web/admin/server` 直接暴露端口。
- `docker-compose.nginx.yaml`：Nginx 统一入口部署，对外暴露 `80/443`，内部转发到各服务。

## 目录说明

- `nginx/conf.d/default.conf`：站点反向代理配置
- `nginx/ssl/`：证书目录，放置证书和私钥文件
- `script/deploy-nginx.sh`：Nginx 入口模式一键部署脚本
- `.env.nginx.example`：Nginx 部署环境变量模板

## 路由规则

- `/` -> `web:3000`
- `/admin` -> `admin:80`
- `/api/*` -> `web:3000`（由 Next.js 继续转发到后端）
- `/static/*` -> `server:8080`

## HTTPS 与安全配置

- `80` 会自动跳转到 `443`
- `443` 启用 TLS，并从 `nginx/ssl/` 读取证书
- 已启用基础安全响应头
- 已启用连接数限制与请求速率限制：
  - 单 IP 最大连接数：`30`
  - 单 IP 请求速率：`10 req/s`
  - 突发请求缓冲：`30`

## 使用方式

1. 复制环境变量模板：
   - `cp .env.nginx.example .env.nginx`
2. 修改 `.env.nginx`：
   - `WEBSITE_BASE_HOST`
   - `WEBSITE_ADMIN_HOST`
   - `WEBSITE_SERVER_HOST`
   - `SERVER_NAME`
   - `SSL_CERT_FILE`
   - `SSL_KEY_FILE`
3. 把证书文件放到 `nginx/ssl/`
4. 执行脚本启动：
   - `bash script/deploy-nginx.sh`

## 证书文件示例

如果你使用默认模板，目录结构类似：

- `nginx/ssl/fullchain.pem`
- `nginx/ssl/privkey.pem`
