package app

import (
	"easy_pwd/global"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
)

func initAppView() fyne.Window {
	// 创建一个新的应用程序实例
	app := app.NewWithID("com.personal.EasyPwd")
	// 获取应用程序的元数据
	appName := app.Metadata().Name
	if appName == "" {
		appName = global.DefProjectName
	}
	// 创建一个新的窗口
	window := app.NewWindow(appName)

	// 设置窗口的图标
	window.SetIcon(theme.ComputerIcon())

	// 设置窗口的初始大小
	window.Resize(fyne.NewSize(960, 720))

	// 设置居中显示
	window.CenterOnScreen()

	// 创建一个多标签页式的容器
	tabs := createTabs(window, app)

	// 设置标签页在上方
	tabs.SetTabLocation(container.TabLocationTop)

	// 将多标签页式的容器添加到窗口中
	window.SetContent(tabs)

	// 设置窗口图标
	window.SetIcon(fyne.NewStaticResource("icon", global.Icon))

	return window
}

func createTabs(w fyne.Window, a fyne.App) *container.AppTabs {
	// 创建密码管理tab容器
	passwordTab := container.NewTabItem("密码管理", createPasswordView(w, a))
	passwordTab.Icon = theme.StorageIcon()

	// 创建密码重置tab容器
	dataTab := container.NewTabItem("数据管理", createDataView(w, a))
	dataTab.Icon = theme.SettingsIcon()

	// 创建版本tab容器
	versionTab := container.NewTabItem("项目信息", createVersionView(w, a))
	versionTab.Icon = theme.InfoIcon()

	// 组装tabs
	tabs := container.NewAppTabs(passwordTab, dataTab, versionTab)
	tabs.SelectIndex(0)
	return tabs
}
