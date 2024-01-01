FROM node:18.18.0-alpine

ENV APP_HOME /fnote/web

WORKDIR "$APP_HOME"
COPY . ./

# install
RUN npm config set registry http://registry.npmmirror.com
RUN npm install

# build
RUN npm run build -- --API_BASE_URL=http://fnote-server:8080

EXPOSE 3000
CMD ["npm", "run", "start"]