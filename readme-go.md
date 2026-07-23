后端go

一、部署项目
1、打包命令
在所在的项目目录下执行
cd D:\GolandProjects
go build -o my_task_app.exe main.go

2、运行程序
双击 my_task_app.exe 或在命令行中运行它。
./my_task_app.exe

3、在windows中注册成服务(开机自启)

二、golandIED工具使用
1、更换依赖所在目录
在所在根目录下执行,默认放在什么位置:
go env GOMODCACHE 

在 Terminal 下执行下面5条命令:
New-Item -ItemType Directory -Path "D:\GolandProjects\go\modcache" -Force
[System.Environment]::SetEnvironmentVariable("GOMODCACHE", "D:\GolandProjects\go\modcache", "User")
go env GOMODCACHE
go clean -modcache
go mod tidy

三、大型项目标准目录结构
在欧美大厂和大型开源项目中广泛采用的标准目录结构

my-project/
├── cmd/                 # 主程序入口目录（支持多个服务）
│   ├── api/             # 核心 Web API 服务
│   │   └── main.go
│   ├── worker/          # 后台任务/消息队列消费者
│   │   └── main.go
│   └── cli/             # 命令行工具（如有）
│       └── main.go
├── internal/            # 内部私有逻辑（核心业务，外部项目无法导入）
│   ├── user/            # 按领域划分模块（DDD思想）
│   │   ├── handler.go   # 请求处理层
│   │   ├── service.go   # 业务逻辑层
│   │   ├── repo.go      # 数据访问层
│   │   └── model.go     # 领域模型
│   ├── order/           # 订单模块
│   └── middleware/      # 内部中间件
├── pkg/                 # 可被外部项目导入的公共库
│   ├── logger/          # 日志封装
│   ├── auth/            # 鉴权工具
│   └── config/          # 配置加载
├── api/                 # 接口定义文件
│   ├── proto/           # gRPC Protobuf 定义
│   └── openapi/         # Swagger/OpenAPI 规范
├── configs/             # 配置文件（yaml, json, env）
├── scripts/             # 构建、部署、数据库迁移脚本
├── test/                # 测试辅助工具或集成测试
├── docs/                # 项目文档
├── go.mod
├── go.sum
└── README.md

