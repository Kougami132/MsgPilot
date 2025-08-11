# 使用轻量级的Alpine Linux作为基础镜像
FROM alpine:3.22

# 设置工作目录
WORKDIR /app

# 安装必要的运行时依赖
RUN apk --no-cache add ca-certificates tzdata

# 设置时区为上海时区
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone

# 复制编译好的Linux可执行文件
COPY dist/msgpilot-linux-amd64 /app/msgpilot

# 复制静态文件
COPY dist/static /app/static

# 创建数据目录并设置权限
RUN mkdir -p /app/data && chmod 755 /app/data

# 设置可执行权限
RUN chmod +x /app/msgpilot

# 暴露端口（根据你的应用配置调整）
EXPOSE 8080

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# 启动应用
CMD ["./msgpilot"]
