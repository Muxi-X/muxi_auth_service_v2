#------ 第一阶段：构建 ----------
FROM golang:1.23-alpine AS builder

# 开启 go module
ENV CGO_ENABLED=0 \
    GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct

WORKDIR /app

COPY . .

RUN go build -o main .


# ---------- 第二阶段：运行 ----------
FROM alpine:latest

WORKDIR /app

# 只复制编译好的二进制
COPY --from=builder /app/main .

COPY conf/config.yaml ./conf/config.yaml

EXPOSE 8083 25 465 587

CMD ["./main", "-c", "conf/config.yaml"]