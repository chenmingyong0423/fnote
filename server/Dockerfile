FROM golang:alpine AS builder

ENV APP_HOME /fnote/server
WORKDIR "$APP_HOME"

COPY . ./

RUN go env -w GOPROXY=https://goproxy.cn,direct && go mod download

RUN go build

FROM alpine:latest

COPY --from=builder /fnote/server/server /fnote/server
COPY --from=builder /fnote/server/config/fnote.prod.yaml /fnote/config/fnote.prod.yaml

WORKDIR /fnote

EXPOSE 8080

ENTRYPOINT ["./server", "--port", ":8080", "--config", "./config/fnote.docker.yaml"]