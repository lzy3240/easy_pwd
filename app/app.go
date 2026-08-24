package app

import (
	"easy_pwd/global"
	"easy_pwd/pkg/db"
	"easy_pwd/pkg/log"
	"easy_pwd/pkg/utils"
	"runtime"
)

func Run() {
	defer func() {
		if r := recover(); r != nil {
			log.Instance().Error("运行失败", log.SetValue("error", r))
			// 获取调用栈信息
			// 2. 获取完整的堆栈信息
			buf := make([]byte, 4096)
			n := runtime.Stack(buf, false) // false: 只获取当前 goroutine 的堆栈
			stackTrace := string(buf[:n])
			//fmt.Println(stackTrace)
			log.Instance().Error("运行失败", log.SetValue("stackTrace", stackTrace))
		}
		// 关闭数据库连接
		db.CloseDB()
	}()

	// 初始化数据库
	checkTable()
	initData()
	// 启动界面
	w := initAppView()
	// 运行应用程序
	w.ShowAndRun()
}

func checkTable() {
	DB := db.Instance()
	// 检查表是否存在, 不存在则创建
	if DB.Migrator().HasTable(&global.PassWord{}) == false {
		if err := DB.Debug().AutoMigrate(&global.PassWord{}); err != nil {
			log.Instance().Error("初始化数据表失败: " + err.Error())
		}
	}
	if DB.Migrator().HasTable(&global.Setting{}) == false {
		if err := DB.Debug().AutoMigrate(&global.Setting{}); err != nil {
			log.Instance().Error("初始化数据表失败: " + err.Error())
		}
	}
}

func initData() {
	var setting global.Setting
	DB := db.Instance()
	err := DB.Model(&global.Setting{}).Where("key = ?", "PrivateKey").First(&setting).Error
	if err != nil {
		log.Instance().Error("获取密钥数据失败: " + err.Error())
	}
	if setting.ID == 0 || setting.Value == "" {
		log.Instance().Error("密钥数据不存在, 将初始化密钥数据")
		// 生成密钥数据
		publicKeyData, privateKeyData, err1 := utils.GetRSAKeys()
		if err1 != nil {
			log.Instance().Error("初始化密钥失败: " + err.Error())
			return
		}
		// 保存密钥数据
		DB.Create(&global.Setting{Key: "PrivateKey", Value: privateKeyData})
		DB.Create(&global.Setting{Key: "PublicKey", Value: publicKeyData})
		log.Instance().Info("初始化密钥成功")
	}

	var setting1 global.Setting
	err = DB.Model(&global.Setting{}).Where("key = ?", "SuperPwdStatus").First(&setting1).Error
	if err != nil {
		log.Instance().Error("获取超级密码状态失败: " + err.Error())
	}
	if setting1.ID == 0 || setting1.Value == "" {
		log.Instance().Error("超级密码启用配置不存在, 将初始化超级密码启用配置")
		// 保存超级密码配置(是否启用)
		DB.Create(&global.Setting{Key: "SuperPwdStatus", Value: "false"})
		log.Instance().Info("初始化超级密码启用配置成功")
	}
}
