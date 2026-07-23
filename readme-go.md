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

