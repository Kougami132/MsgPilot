# MsgPilot

MsgPilot是一个基于Go语言开发的消息处理系统，提供统一的消息管理和转发服务。

## 特性

- 基于Gin框架的RESTful API
- JWT身份认证
- SQLite数据库存储
- Swagger API文档
- 消息通道管理
- 适配器系统

## 技术栈

- Go 1.23.2
- Gin Web框架
- GORM
- SQLite
- JWT
- Swagger

## 快速开始

1. 克隆项目
```bash
git clone https://github.com/kougami132/MsgPilot.git
cd MsgPilot
```

2. 安装依赖
```bash
go mod download
```

3. 配置环境变量
复制`.env.example`到`.env`并根据需要修改配置：
```env
APP_ENV=development
PORT=8080
CONTEXT_TIMEOUT=10
ACCESS_TOKEN_EXPIRY_HOUR=720
ACCESS_TOKEN_SECRET=secret
FRONTEND_PATH=./frontend/dist
```

4. 运行服务
```bash
go run main.go
```

服务将在`http://localhost:8080`启动

## Docker部署

### 使用Docker Compose（推荐）

1. 确保已安装Docker和Docker Compose

2. 编辑docker-compose.yml文件
```
services:
    msgpilot:
        image: kougami132/msgpilot:latest
        container_name: msgpilot
        restart: unless-stopped
        network_mode: bridge
        ports:
            - "8080:8080"
        volumes:
            - ./data:/app/data
        environment:
            - ACCESS_TOKEN_SECRET=secret
```

3. 启动服务
```bash
docker compose up -d
```

4. 查看服务状态
```bash
docker compose ps
docker compose logs -f msgpilot
```

5. 停止服务
```bash
docker compose down
```

### 使用Docker

```bash
docker run -d \
  --name msgpilot \
  -p 8080:8080 \
  -e ACCESS_TOKEN_SECRET=secret \
  -v ./data:/app/data \
  kougami132/msgpilot:latest
```

## API文档

启动服务后访问`http://localhost:8080/api/swagger/index.html`查看API文档

## 项目结构

```
.
├── api/            # API处理器和中间件
├── bootstrap/      # 应用初始化和依赖注入
├── config/         # 配置文件
├── docs/          # API文档
├── frontend/      # 前端代码
├── internal/       # 内部包
├── models/         # 数据模型
├── route/         # 路由定义
├── test/          # 测试文件
└── main.go        # 程序入口
```

## License

本项目基于GPL-2.0许可证开源