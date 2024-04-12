FROM node:18.18.0-alpine AS builder

ENV APP_HOME /fnote/admin

WORKDIR "$APP_HOME"
COPY . ./

# install
RUN npm config set registry http://registry.npmmirror.com
RUN npm install
RUN npm run build:docker

FROM nginx:latest

COPY --from=builder /fnote/admin/dist /usr/share/nginx/html

COPY nginx.conf /etc/nginx/conf.d/default.conf


EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]