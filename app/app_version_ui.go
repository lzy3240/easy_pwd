package app

import (
	"easy_pwd/global"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"net/url"
)

func createVersionView(w fyne.Window, a fyne.App) *fyne.Container {
	nameCard := container.NewHBox(
		widget.NewLabelWithStyle("工具名称", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewSeparator(),
		widget.NewLabelWithStyle(global.DefProjectName, fyne.TextAlignCenter, fyne.TextStyle{Bold: false}),
	)
	nameCard = container.NewPadded(nameCard)

	versionCard := container.NewHBox(
		widget.NewLabelWithStyle("版本信息", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewSeparator(),
		widget.NewLabelWithStyle(global.Version, fyne.TextAlignCenter, fyne.TextStyle{Bold: false}),
	)
	versionCard = container.NewPadded(versionCard)

	//authorCard := container.NewHBox(
	//	widget.NewLabelWithStyle("软件作者", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
	//	widget.NewSeparator(),
	//	widget.NewLabelWithStyle(global.Author, fyne.TextAlignCenter, fyne.TextStyle{Bold: false}),
	//)
	//authorCard = container.NewPadded(authorCard)
	//
	//contactCard := container.NewHBox(
	//	widget.NewLabelWithStyle("联系方式", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
	//	widget.NewSeparator(),
	//	widget.NewLabelWithStyle(global.Contact, fyne.TextAlignCenter, fyne.TextStyle{Bold: false}),
	//)
	//contactCard = container.NewPadded(contactCard)

	// 创建项目链接按钮
	urlButton := widget.NewButtonWithIcon("项目地址", theme.HelpIcon(), func() {
		parsedURL, err := url.Parse(global.ProjectURL)
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		if err = a.OpenURL(parsedURL); err != nil {
			dialog.ShowError(err, w)
			return
		}
	})
	urlButton.Importance = widget.HighImportance
	urlCard := container.NewPadded(container.NewVBox(urlButton))

	return container.NewCenter(
		container.NewVBox(
			widget.NewLabelWithStyle("欢迎使用", fyne.TextAlignCenter, fyne.TextStyle{
				Bold:      true,
				Underline: true,
			}),
			layout.NewSpacer(),
			nameCard,
			layout.NewSpacer(),
			versionCard,
			layout.NewSpacer(),
			//authorCard,
			//layout.NewSpacer(),
			//contactCard,
			urlCard,
			layout.NewSpacer(),
		),
	)
}
