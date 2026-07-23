package config

import (
	"GolandProjects/model"
	"fmt" // 用于格式化 DSN
	// 1. PostgreSQL 驱动
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	// 2. 修改连接信息
	host := "127.0.0.1"
	user := "postgres"
	password := "Admin@9000"
	dbname := "postgres"
	port := 5432 // PostgreSQL 默认端口是 5432

	// 3. 使用 postgres 驱动和新的 DSN 格式
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
		host, user, password, dbname, port)

	var err error
	// 注意这里从 mysql.Open 变成了 postgres.Open
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("数据库连接失败: " + err.Error())
	}
	// 自动建表
	DB.AutoMigrate(&model.TaskRecord{})
	DB.AutoMigrate(&model.TaskNewRecord{})
}

//package config
//
//import (
//	"context"
//	"fmt"
//	"log"
//	"time"
//
//	// GORM 核心库
//	"gorm.io/gorm"
//	// PostgreSQL 驱动
//	"gorm.io/driver/postgres"
//	// MySQL 驱动
//	"github.com/redis/go-redis/v9"
//	"gorm.io/driver/mysql"
//)
//
//// DBConfig 用于存放所有数据库连接实例
//type DBConfig struct {
//	MasterPG    *gorm.DB        // PostgreSQL 主库
//	SlaveMySQL  *gorm.DB        // MySQL 从库
//	RedisClient *redis.Client   // Redis 客户端
//	Ctx         context.Context // Redis 操作所需的上下文
//}
//
//var AppConfig *DBConfig
//
//// InitDatabases 初始化所有数据库连接
//func InitDatabases() {
//	AppConfig = &DBConfig{
//		Ctx: context.Background(),
//	}
//
//	// 1. 初始化 PostgreSQL 主库
//	initPostgreSQL()
//
//	// 2. 初始化 MySQL 从库
//	initMySQL()
//
//	// 3. 初始化 Redis
//	initRedis()
//}
//
//// initPostgreSQL 连接 PostgreSQL 主库
//func initPostgreSQL() {
//	// 请替换为你的 PostgreSQL 实际连接信息
//	dsn := "host=127.0.0.1 user=postgres password=Admin@9000 dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Shanghai"
//
//	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
//	if err != nil {
//		log.Fatalf("❌ 连接 PostgreSQL 主库失败: %v", err)
//	}
//
//	// 获取底层的 sql.DB 以配置连接池
//	sqlDB, err := db.DB()
//	if err != nil {
//		log.Fatalf("❌ 获取 PostgreSQL sql.DB 实例失败: %v", err)
//	}
//	sqlDB.SetMaxIdleConns(10)
//	sqlDB.SetMaxOpenConns(100)
//	sqlDB.SetConnMaxLifetime(time.Hour)
//
//	AppConfig.MasterPG = db
//	fmt.Println("✅ PostgreSQL 主库连接成功")
//}
//
//// initMySQL 连接 MySQL 从库
//func initMySQL() {
//	// 请替换为你的 MySQL 实际连接信息
//	// 注意：这里添加了 parseTime=true 以正确处理时间类型
//	dsn := "user:password@tcp(127.0.0.1:3306)/your_db?charset=utf8mb4&parseTime=True&loc=Local"
//
//	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
//	if err != nil {
//		log.Fatalf("❌ 连接 MySQL 从库失败: %v", err)
//	}
//
//	// 获取底层的 sql.DB 以配置连接池
//	sqlDB, err := db.DB()
//	if err != nil {
//		log.Fatalf("❌ 获取 MySQL sql.DB 实例失败: %v", err)
//	}
//	sqlDB.SetMaxIdleConns(10)
//	sqlDB.SetMaxOpenConns(100)
//	sqlDB.SetConnMaxLifetime(time.Hour)
//
//	AppConfig.SlaveMySQL = db
//	fmt.Println("✅ MySQL 从库连接成功")
//}
//
//// initRedis 连接 Redis
//func initRedis() {
//	// 请替换为你的 Redis 实际连接信息
//	rdb := redis.NewClient(&redis.Options{
//		Addr:     "localhost:6379",
//		Password: "", // 没有密码则留空
//		DB:       0,  // 使用默认 DB
//	})
//
//	// 测试连接
//	_, err := rdb.Ping(AppConfig.Ctx).Result()
//	if err != nil {
//		log.Fatalf("❌ 连接 Redis 失败: %v", err)
//	}
//
//	AppConfig.RedisClient = rdb
//	fmt.Println("✅ Redis 连接成功")
//}
