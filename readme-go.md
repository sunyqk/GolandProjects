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
1、在欧美大厂和大型开源项目中广泛采用的标准目录结构

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

2、核心设计原则与规范
cmd/ 作为多服务入口
大型后端通常不是单一程序，而是包含 API 网关、后台 Worker、定时任务等多个二进制文件。cmd/ 下每个子目录对应一个 main.go，负责组装依赖并启动服务。
internal/ 保证代码私有性
这是 Go 语言的核心访问控制机制。internal/ 下的代码只能被其父目录及其子目录导入。大型项目中，核心业务逻辑（如 user、order）必须放在这里，防止被其他微服务或第三方项目意外引用。
按“领域”而非“技术”划分模块
在 internal/ 下，欧美团队极力推崇领域驱动设计（DDD）。不要建全局的 handler/、service/、model/ 文件夹，而是以业务实体为核心（如 internal/user/），将处理、逻辑、数据访问集中在一起，降低跨模块耦合。
pkg/ 存放公共库
只有当你确定某些代码（如加密工具、通用校验器）需要被其他独立仓库导入时，才放入 pkg/。现代 Go 项目已不再滥用 pkg/，如果代码仅在内部使用，请一律放在 internal/。
拒绝 util/ 和 common/
大型项目严禁建立名为 util/、common/ 或 helper/ 的“万能包”。这些文件夹最终会变成毫无逻辑的代码垃圾场。工具函数应放在语义相关的包中（如 pkg/crypto/），或直接作为私有方法写在业务包内部。

3、生产环境部署目录（Linux 服务器）
当你的项目打包部署到 Linux 服务器（如 Ubuntu）时，欧美运维团队通常遵循 Linux FHS 规范，将应用部署在 /opt 目录下，而不是源码目录：

/opt/my-project/
├── bin/                 # 编译后的二进制文件（绝不上传 .go 源码）
├── public/              # 静态资源（由 Nginx 直接服务）
├── media/               # 用户上传的文件（需独立挂载写权限）
└── config.yaml          # 运行时配置（与代码分离）

关键原则：二进制与源码必须分离；日志输出到 /var/log/my-project/；可写数据（如缓存、上传文件）与代码物理隔离。
按照这套结构重构你的项目，不仅能满足大型项目的扩展需求，还能让你的代码看起来非常符合国际一线大厂的工程规范。