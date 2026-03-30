FROM node:18.18.0-alpine AS builder

ENV APP_HOME /fnote/admin
WORKDIR "$APP_HOME"
COPY . ./

# build args from docker-compose
ARG VITE_BASE_HOST
ARG VITE_API_HOST

# expose to Vite build process
ENV VITE_BASE_HOST=${VITE_BASE_HOST}
ENV VITE_API_HOST=${VITE_API_HOST}

# install
RUN npm config set registry http://registry.npmmirror.com
RUN npm install
RUN npm run build:docker

FROM nginx:latest

COPY --from=builder /fnote/admin/dist /usr/share/nginx/html

COPY default.conf /etc/nginx/conf.d/default.conf


EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]