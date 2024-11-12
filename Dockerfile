# 第一阶段：构建 Go 应用
FROM golang:alpine AS builder

# 设置 Go 代理为七牛云的代理
ENV GOPROXY=https://goproxy.cn,direct

# 安装需要的依赖
RUN apk update && apk add --no-cache git

# 设置工作目录为 /app
WORKDIR /app

# 复制源代码
COPY . ./

COPY ./file ./file

# 执行依赖更新并构建二进制文件
RUN go mod tidy && go build -o main .


# 第二阶段：构建最终镜像
FROM alpine

# 设置工作目录为 /app
WORKDIR /app


# 从 builder 镜像复制编译好的二进制文件
COPY --from=builder /app/main .


# 复制图片文件到最终镜像
COPY --from=builder /app/file ./file

# 暴露端口
EXPOSE 8000

# 启动应用程序
CMD ["./main"]
