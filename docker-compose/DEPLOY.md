# 部署说明

本文档说明 Nginx 统一入口模式的部署步骤。

## 准备工作

1. 在服务器安装 Docker 和 Docker Compose。
2. 将域名解析到服务器 IP，并开放 `80`、`443` 端口。
3. 准备 HTTPS 证书和私钥，例如：
   - `fullchain.pem`
   - `privkey.pem`

## 配置环境变量

进入 `docker-compose` 目录，复制环境变量模板：

```bash
cp .env.nginx.example .env.nginx
```

修改 `.env.nginx`：

```env
WEBSITE_BASE_HOST=https://你的域名
WEBSITE_ADMIN_HOST=https://你的域名/admin
WEBSITE_SERVER_HOST=https://你的域名
SERVER_NAME=你的域名
SSL_CERT_FILE=fullchain.pem
SSL_KEY_FILE=privkey.pem
```

## 放置证书

将证书和私钥放到：

```text
docker-compose/nginx/ssl/
```

文件名需要和 `.env.nginx` 中的 `SSL_CERT_FILE`、`SSL_KEY_FILE` 保持一致。

## 启动部署

在 `docker-compose` 目录执行：

```bash
bash script/deploy-nginx.sh
```

脚本会检查必填环境变量和证书文件，然后构建并启动服务。

## 检查服务

```bash
docker compose --env-file .env.nginx -f docker-compose.nginx.yaml ps
docker compose --env-file .env.nginx -f docker-compose.nginx.yaml logs -f nginx
docker compose --env-file .env.nginx -f docker-compose.nginx.yaml logs -f server
```

访问：

```text
https://你的域名
https://你的域名/admin
```

