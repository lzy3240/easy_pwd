# easy_pwd 密码管理工具

## 📋 项目简介
**easy_pwd** 是一款跨平台密码管理工具，采用 Go 语言开发，使用 Fyne 框架构建的图形界面工具。旨在为用户提供简单、高效的密码存储和管理解决方案，支持本地 SQLite 数据库存储，确保数据安全性和隐私保护。

## 📦 安装指南
### 环境要求

| 组件   | 最低版本   | 推荐版本    | 说明        |
| ---- | ------ | ------- | --------- |
| Go   | 1.25+  | 1.25.0+ | Go 编程语言环境 |
| Fyne | 2.8.0+ | 2.8.0+  | 用于构建图形界面  |
| Git  | 2.0+   | 最新版     | 用于克隆代码仓库  |

### 下载源码
#### 1. 克隆仓库

```bash
git clone https://gitee.com/lzy3240/easy_pwd.git
cd easy_pwd
```

#### 2. 安装依赖

```bash
go mod download
```

#### 3. 运行应用

```bash
go run main.go
```

### 程序构建
**强烈建议：
使用常规构建方法, 构建速度更快、更简单, 生成的可执行文件更小**

#### 1.常规构建

```bash
go build -ldflags "-s -w -H=windowsgui"
```

#### 2.常规构建带Icon
- 进入到项目的目录
- 在项目入口目录内，制作ico图片,如main.ico, 不可直接使用png图片, 必须转换成ico格式
- 在项目入口目录内，创建一个空白文本文件,命名main.rc,内容输入: IDI_ICON1 ICON "main.ico"
- 在项目入口目录内，执行命令: windres -o main.syso main.rc ,此时生成了一个main.syso
- go build编译即可, 生成带图标的执行文件
 
```bash
执行 go build -ldflags "-s -w -H=windowsgui"
```

#### 3.Fyne构建
```bash
fyne install -icon ./icon.png                #检测当前系统打包程序，不指定png默认是Icon.png
fyne package -os darwin -icon ./icon.png     #macOS系统,创建myapp.app
fyne package -os linux -icon ./icon.png      #linux系统，创建myapp.tar.gz
fyne package -os windows -icon ./icon.png     #window系统，创建myapp.exe
```


## 🛠️ 技术特性
- ✅ **SQLite 数据库**：轻量级本地数据库，无CGO
- ✅ **数据加密**：采用RSA加密存储，保护用户隐私
- ✅ **自动时间戳**：自动记录创建和更新时间
- ✅ **日志系统**：完善的日志记录，便于问题排查
- ✅ **跨平台构建**：支持批量构建多个平台版本
