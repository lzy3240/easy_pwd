package db

import (
	"easy_pwd/global"
	"easy_pwd/pkg/log"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	defaultLog "log"
)

// sqlite放弃CGO驱动, 采用纯go驱动, 避免CGO编译问题; 但编译后体积增大, 性能存在损失.

var (
	db  *gorm.DB
	err error
)

func Instance() *gorm.DB {
	if db == nil {
		InitConn()
	}
	return db
}

func InitConn() {
	gormSqlite()
}

func CloseDB() {
	if db != nil {
		sqlDB, _ := db.DB()
		if sqlDB.Close() != nil {
			log.Instance().Error("关闭数据库异常: " + sqlDB.Close().Error())
		}
	}
}

//func gormMysql() {
//	m := config.Instance().DB
//	dsn := m.DBUser + ":" + m.DBPwd + "@tcp(" + m.DBHost + ":" + m.DBPort + ")/" + m.DBName + "?charset=utf8mb4&parseTime=True&loc=Local"
//	db, err = gorm.Open(mysql.Open(dsn),
//		&gorm.Config{
//			// 启用更新时间戳功能
//			NowFunc: func() time.Time {
//				return time.Now().Local()
//			},
//			// 启用单数表和前缀
//			NamingStrategy: schema.NamingStrategy{
//				//TablePrefix:   "pre_", // 表前缀
//				SingularTable: true, // 禁用表名复数
//			},
//			// 启用数据库日志
//			Logger: logger.New(defaultLog.New(os.Stdout, "\r\n", defaultLog.LstdFlags),
//				logger.Config{
//					SlowThreshold: time.Second,
//					LogLevel:      logger.Info,
//					Colorful:      true,
//				}),
//		})
//	if err != nil {
//		log.Instance().Info("MySQL启动异常: ", err.Error())
//		os.Exit(0)
//	}
//
//	sqlDB, _ := db.DB()
//	sqlDB.SetMaxIdleConns(100)
//	sqlDB.SetMaxOpenConns(300)
//	sqlDB.SetConnMaxLifetime(time.Hour)
//}

func gormSqlite() {
	if checkNotExist(global.DefRuntimeDir) {
		if err = os.MkdirAll(global.DefRuntimeDir, os.ModePerm); err != nil {
			log.Instance().Error("创建数据库目录失败: " + err.Error())
		}
	}
	dbFile := path.Join(global.DefRuntimeDir, fmt.Sprintf("%s.db", global.DBFileName))
	if checkNotExist(dbFile) {
		if err = createDB(dbFile); err != nil {
			log.Instance().Error("创建数据库文件失败: " + err.Error())
		}
	}
	db, err = gorm.Open(sqlite.Open(dbFile),
		&gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true, //禁用外键
			SkipDefaultTransaction:                   true, //禁用事务
			PrepareStmt:                              true, // 预编译语句
			// 启用更新时间戳功能
			NowFunc: func() time.Time {
				return time.Now().Local()
			},
			// 启用单数表和前缀
			NamingStrategy: schema.NamingStrategy{
				//TablePrefix:   "pre_", // 表前缀
				SingularTable: true, // 禁用表名复数
			},
			// 启用数据库日志
			Logger: logger.New(defaultLog.New(os.Stdout, "\r\n", defaultLog.LstdFlags),
				logger.Config{
					SlowThreshold: time.Second,
					LogLevel:      logger.Info,
					Colorful:      true,
				}),
		})
	if err != nil {
		log.Instance().Error("打开数据库失败: " + err.Error())
		os.Exit(0)
	}
	//db.Exec("PRAGMA journal_mode=WAL;") // 开启wal模式, 提升并发性能
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(100)
	sqlDB.SetMaxOpenConns(300)
	sqlDB.SetConnMaxLifetime(time.Hour)
}

func createDB(path string) error {
	fp, err := os.Create(path) // 如果文件已存在，会将文件清空。
	if err != nil {
		return err
	}
	defer fp.Close() //关闭文件，释放资源。
	return nil
}

func checkNotExist(src string) bool {
	_, err = os.Stat(src)
	return os.IsNotExist(err)
}
