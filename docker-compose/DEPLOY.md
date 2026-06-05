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

MongoDB 默认使用内置账号密码。需要自定义时，可以在 `.env.nginx` 里补充：

```env
MONGO_ROOT_USERNAME=fnote
MONGO_ROOT_PASSWORD=fnote
MONGO_DATABASE=fnote
MONGO_USERNAME=fnote-user
MONGO_PASSWORD=你的密码
```

如果 `data/mongo` 已经初始化过，修改这些变量不会自动修改已有 MongoDB 用户，需要手动改密码或清空数据后重新初始化。

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

## 单独重建服务

只需要重新构建某个服务时，在 `docker-compose` 目录执行：

```bash
bash script/rebuild-nginx-service.sh web
bash script/rebuild-nginx-service.sh admin
bash script/rebuild-nginx-service.sh server
```

脚本只会重新构建并启动指定服务，不会重建依赖服务。Nginx 使用官方镜像，没有构建步骤；如果只改了 Nginx 配置，执行：

```bash
docker compose --env-file .env.nginx -f docker-compose.nginx.yaml up -d --force-recreate nginx
```

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

## 上传限制

Nginx 默认请求体限制为 `20m`。后台文件上传接口允许 `100m`，备份恢复接口允许 `512m`：

```text
/admin-api/files/upload
/admin-api/recovery
```

## 数据目录

Nginx 部署的数据会挂载到 `docker-compose/data/`：

```text
data/mongo
data/logs
data/static
```

如果旧版本已经使用 `/tmp/fnote` 保存过数据，部署新版前先迁移数据：

```bash
mkdir -p data
cp -a /tmp/fnote/mongo_data data/mongo
cp -a /tmp/fnote/logs data/logs
cp -a /tmp/fnote/static data/static
```

## 停止服务

在 `docker-compose` 目录执行：

```bash
bash script/stop-nginx.sh
```

脚本会停止并移除 Nginx 部署相关容器和网络，不会删除 `data/` 下的数据目录。

## 重启服务

在 `docker-compose` 目录执行：

```bash
bash script/restart-nginx.sh
```

脚本会重启当前 Nginx 部署相关容器，不会重新构建镜像，也不会删除数据目录。
