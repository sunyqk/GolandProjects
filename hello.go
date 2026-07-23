package main

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// 1. 定义接口：实现多态
type Task interface {
	Execute() (string, error)
}

// 2. 定义结构体：实现接口
type DataFetchTask struct {
	ID int
}

func (t *DataFetchTask) Execute() (string, error) {
	// 模拟耗时操作
	duration := time.Duration(rand.Intn(3)) * time.Second
	time.Sleep(duration)

	// 模拟随机失败
	if rand.Intn(10) < 3 {
		return "", fmt.Errorf("task %d 执行超时", t.ID)
	}
	return fmt.Sprintf("Task-%d 成功获取数据 (耗时 %v)", t.ID, duration), nil
}

// 3. 核心处理函数：结合并发、Channel 和 WaitGroup
func ProcessTasks(tasks []Task) []string {
	var (
		wg      sync.WaitGroup
		mu      sync.Mutex // 互斥锁，保护并发写入
		results []string
	)

	// 使用无缓冲或有缓冲的 channel 都可以，这里用 WaitGroup + Mutex 演示
	for _, task := range tasks {
		wg.Add(1)
		// 闭包捕获 task，注意 Go 1.22 之前的循环变量陷阱（这里显式传参更安全）
		go func(t Task) {
			defer wg.Done() // 4. defer 确保计数器减一

			result, err := t.Execute()

			// 5. 错误处理
			if err != nil {
				fmt.Printf("[警告] %v\n", err)
				return
			}

			// 并发安全地写入结果
			mu.Lock()
			results = append(results, result)
			mu.Unlock()
		}(task)
	}

	wg.Wait() // 阻塞，直到所有 goroutine 执行完毕
	return results
}

func main() {
	fmt.Println("=== Go 复杂语法演示开始 ===")

	// 构造任务列表
	tasks := make([]Task, 5)
	for i := 0; i < 5; i++ {
		tasks[i] = &DataFetchTask{ID: i + 1}
	}

	// 执行并发处理
	start := time.Now()
	results := ProcessTasks(tasks)
	elapsed := time.Since(start)

	// 输出结果
	fmt.Printf("\n处理完成，共成功 %d 个任务，总耗时: %v\n", len(results), elapsed)
	for _, res := range results {
		fmt.Println(" ->", res)
	}

	// 演示自定义错误包装
	err := simulateDeepError()
	if err != nil {
		fmt.Printf("\n[深层错误捕获]: %v\n", err)
	}
}

// 演示 errors 的包装与解包
func simulateDeepError() error {
	baseErr := errors.New("数据库连接失败")
	return fmt.Errorf("业务层处理失败: %w", baseErr)
}
