FROM golang:alpine AS builder

ENV APP_HOME /fnote/backend
WORKDIR "$APP_HOME"

COPY . ./

RUN go env -w GOPROXY=https://goproxy.cn,direct && go mod download

RUN go build

FROM alpine:latest

COPY --from=builder /fnote/backend/backend /fnote/backend
COPY --from=builder /fnote/backend/config/fnote.prod.yaml /fnote/config/fnote.prod.yaml

WORKDIR /fnote

EXPOSE 8080

ENTRYPOINT ["./backend", "--port", ":8080", "--config", "./config/fnote.prod.yaml"]