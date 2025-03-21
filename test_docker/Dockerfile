FROM golang:tip-alpine3.21 as builder

WORKDIR /build

# 设置环境变量，优化 Go 构建
ENV CGO_ENABLED=0 \
    GOPROXY=https://goproxy.cn,direct

# 复制项目文件（假设 main.go 在根目录）
COPY main.go .
COPY go.mod .
COPY go.sum .

# 构建可执行文件
RUN go build -o main ./main.go

# 使用更小的 Alpine 作为运行环境
FROM alpine

# 设置时区
ENV TZ Asia/Shanghai

WORKDIR /data/apps/

# 复制构建产物到最终镜像
COPY --from=builder /build/main .

# 运行应用
ENTRYPOINT ["./main"]