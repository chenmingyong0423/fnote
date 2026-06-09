# 安装依赖
FROM node:24-alpine AS deps
WORKDIR /app
COPY package.json package-lock.json* ./
RUN npm install

# 构建阶段
FROM node:24-alpine AS builder
WORKDIR /app
ENV NEXT_TELEMETRY_DISABLED=1

ARG SERVER_HOST
ARG BASE_HOST
ARG ADMIN_HOST
ARG NEXT_PUBLIC_SERVER_HOST
ENV SERVER_HOST=${SERVER_HOST}
ENV BASE_HOST=${BASE_HOST}
ENV ADMIN_HOST=${ADMIN_HOST}
ENV NEXT_PUBLIC_SERVER_HOST=${NEXT_PUBLIC_SERVER_HOST}

COPY --from=deps /app/node_modules ./node_modules
COPY . .
RUN npm run build

# 运行阶段
FROM node:24-alpine AS runner
WORKDIR /app
ENV NODE_ENV=production
ENV NEXT_TELEMETRY_DISABLED=1

USER node

COPY --from=builder /app/.next ./.next
COPY --from=builder /app/package.json ./package.json
COPY --from=builder /app/node_modules ./node_modules

EXPOSE 3000
CMD ["npm", "run", "start"]
