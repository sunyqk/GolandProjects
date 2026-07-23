package job

import (
	"GolandProjects/pkg/logger"
	"fmt"
	"time"

	"go.uber.org/zap"
)

func StartTaskMonitor() {
	ticker := time.NewTicker(1 * time.Minute)
	go func() {
		//for range ticker.C {
		//	// 调用 service 层处理业务，不要直接在定时器里写 SQL
		//	service.CreateTask();
		//}

		// 打印启动信息
		fmt.Println(">>> 任务监控定时器已启动...")

		for range ticker.C {
			// 2. 在循环内打印当前时间，确认它在跑
			fmt.Printf("⏰ 定时器触发 | 当前时间: %s\n", time.Now().Format("2006-01-02 15:04:05"))

			// 这里调用你的业务逻辑
			// service.CreateTask()

			// 测试模拟全局日志配置
			simulateBusinessLogic()
		}
	}()
}

func simulateBusinessLogic() {
	name := "张三"
	// 上述要是输入日志系统中需要用 zap.String() , 或其他数据结构来转换
	logger.Info("测试全局日志系统已经启动..............", zap.String("name", name))
}
