package app

import (
	"easy_pwd/global"
	"easy_pwd/pkg/db"
	"easy_pwd/pkg/log"
	"easy_pwd/pkg/utils"
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/xuri/excelize/v2"
	"image/color"
	"strconv"
	"time"
)

type forcedVariant struct {
	fyne.Theme
	variant fyne.ThemeVariant
}

func (f *forcedVariant) Color(name fyne.ThemeColorName, _ fyne.ThemeVariant) color.Color {
	return f.Theme.Color(name, f.variant)
}

func createDataView(w fyne.Window, a fyne.App) *fyne.Container {
	middleContainer := container.NewVBox(
		createDataContainer(w, a),
	)
	return container.NewBorder(
		nil,
		nil,
		nil,
		nil,
		middleContainer,
	)
}

func createDataContainer(w fyne.Window, a fyne.App) *fyne.Container {
	// 创建密码导出按钮
	excelExportButton := createExcelExportButton(w)
	excelExportButton.Importance = widget.HighImportance
	excelExportButton.SetIcon(theme.DownloadIcon())

	// 创建密码导入按钮
	importExcelButton := createExcelImportButton(w)
	importExcelButton.Importance = widget.HighImportance
	importExcelButton.SetIcon(theme.UploadIcon())

	// 创建超密设置按钮
	editSuperPwdButton := createEditSuperPwdButton(w)
	editSuperPwdButton.Importance = widget.HighImportance
	editSuperPwdButton.SetIcon(theme.MediaReplayIcon())

	// 创建忘记超级密码按钮
	forgetSuperPwdButton := createForgetSuperPwdButton(w)
	forgetSuperPwdButton.Importance = widget.HighImportance
	forgetSuperPwdButton.SetIcon(theme.HistoryIcon())

	// 创建超级密码状态按钮
	statusSuperPwdButton := createStatusSuperPwdButton(w)
	statusSuperPwdButton.Importance = widget.HighImportance
	statusSuperPwdButton.SetIcon(theme.SettingsIcon())

	// 创建密码数据清空按钮
	clearDataButton := createClearDataButton(w)
	clearDataButton.Importance = widget.WarningImportance
	clearDataButton.SetIcon(theme.DeleteIcon())

	// 创建暗黑主题按钮
	darkThemeButton := widget.NewButton("切换暗黑主题", func() {
		a.Settings().SetTheme(&forcedVariant{Theme: theme.DefaultTheme(), variant: theme.VariantDark})
	})
	darkThemeButton.SetIcon(theme.ViewRestoreIcon())
	lightThemeButton := widget.NewButton("切换明亮主题", func() {
		a.Settings().SetTheme(&forcedVariant{Theme: theme.DefaultTheme(), variant: theme.VariantLight})
	})
	lightThemeButton.SetIcon(theme.ViewFullScreenIcon())

	card1 := container.NewHBox(
		widget.NewLabelWithStyle("导出导入", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		//表格导出按钮
		container.NewPadded(excelExportButton),
		//表格导入按钮
		container.NewPadded(importExcelButton),
	)
	exportOpt := container.NewPadded(card1)

	card2 := container.NewHBox(
		widget.NewLabelWithStyle("超级密码", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		//超级密码修改
		container.NewPadded(editSuperPwdButton),
		//忘记超级密码
		container.NewPadded(forgetSuperPwdButton),
	)
	superPwdOpt := container.NewPadded(card2)

	card3 := container.NewHBox(
		widget.NewLabelWithStyle("其他操作", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		//超级密码状态
		container.NewPadded(statusSuperPwdButton),
		//密码数据清空
		container.NewPadded(clearDataButton),
	)
	otherOpt := container.NewPadded(card3)

	card4 := container.NewHBox(
		widget.NewLabelWithStyle("主题切换", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		container.NewPadded(darkThemeButton),
		container.NewPadded(lightThemeButton),
	)
	themeOpt := container.NewPadded(card4)

	return container.NewCenter(
		container.NewVBox(
			widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{}),
			widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{}),
			widget.NewLabelWithStyle("数据管理", fyne.TextAlignCenter, fyne.TextStyle{
				Bold:      true,
				Underline: true,
			}),
			layout.NewSpacer(),
			exportOpt,
			layout.NewSpacer(),
			superPwdOpt,
			layout.NewSpacer(),
			otherOpt,
			layout.NewSpacer(),
			themeOpt,
			layout.NewSpacer(),
		),
	)
}

// --------------------------按钮层--------------------------
func createExcelExportButton(w fyne.Window) *widget.Button {
	return widget.NewButton("密码表格导出", func() {
		// 创建一个文件选择对话框
		saveFile := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
			// 检查错误
			if err != nil {
				dialog.ShowError(err, w)
				return
			}

			// 检查文件写入器
			if writer == nil {
				return
			}
			defer func() { _ = writer.Close() }()
			// 带加载动画的载入过程
			loading := dialog.NewCustomWithoutButtons("载入中...", widget.NewProgressBarInfinite(), w)

			go func() {
				// 显示加载动画
				fyne.Do(func() {
					loading.Show()
				})

				// TODO
				// 查询密码数据
				queryPwdData()
				// 导出密码数据
				err = exportPwdToFile(pwdDataList, writer)
				if err != nil {
					dialog.ShowError(err, w)
					return
				}
				// 隐藏加载动画并显示导入成功信息
				fyne.Do(func() {
					loading.Hide()
					dialog.ShowInformation("导出成功", "成功导出 "+strconv.Itoa(len(pwdDataList))+" 条密码数据", w)
				})
			}()
		}, w)
		// 设置文件选择对话框的默认文件名
		saveFile.SetFileName("password_data.xlsx")

		// 设置文件选择对话框的过滤器
		saveFile.SetFilter(storage.NewExtensionFileFilter([]string{".xlsx"}))

		// 设置文件选择对话框的样式
		saveFile.SetDismissText("取消")
		saveFile.SetConfirmText("保存")

		saveFile.Show()
		saveFile.Resize(fyne.NewSize(800, 600)) // v2.8.0后需要先显示对话框再调整大小

		// 如果查询密码为空
		if len(pwdDataList) == 0 {
			dialog.ShowError(fmt.Errorf("没有记录可以导出"), w)
			return
		}
	})
}

func createExcelImportButton(w fyne.Window) *widget.Button {
	return widget.NewButton("密码表格导入", func() {
		// 创建文件选择器对话框
		openFile := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			// 检查错误
			if err != nil {
				dialog.ShowError(err, w)
				return
			}

			// 检查文件读取器
			if reader == nil {
				return
			}
			defer func() { _ = reader.Close() }()

			// 带加载动画的载入过程
			loading := dialog.NewCustomWithoutButtons("载入中...", widget.NewProgressBarInfinite(), w)

			// 执行导入, 考虑数据量较大，使用goroutine进行异步处理
			go func() {
				// 显示加载动画
				fyne.Do(func() {
					loading.Show()
				})

				// 打开excel导入
				err1, count := importPwdFromFile(reader.URI().Path())
				if err1 != nil {
					dialog.ShowError(errors.New("导入文件失败"+err1.Error()), w)
					return
				}

				// 刷新密码数据
				queryPwdData()

				// 隐藏加载动画并显示导入成功信息
				fyne.Do(func() {
					loading.Hide()
					dialog.ShowInformation("导入成功", "成功导入 "+strconv.Itoa(count)+" 条密码数据", w)
				})
			}()
		}, w)

		// 设置文件选择对话框的过滤器
		openFile.SetFilter(storage.NewExtensionFileFilter([]string{".xlsx"}))

		// 设置文件选择对话框的样式
		openFile.SetDismissText("取消")
		openFile.SetConfirmText("选择")

		openFile.Show()

		// v2.8.0调整, 需要先显示对话框再调整大小
		openFile.Resize(fyne.NewSize(800, 600))
	})
}

func createEditSuperPwdButton(w fyne.Window) *widget.Button {
	return widget.NewButton("超级密码修改", func() {
		oldSuperPwdEntry := widget.NewPasswordEntry()
		oldSuperPwdEntry.SetPlaceHolder("请输入原超级密码")

		newSuperPwdEntry := widget.NewPasswordEntry()
		newSuperPwdEntry.SetPlaceHolder("请输入新超级密码")
		newSuperPwdEntry2 := widget.NewPasswordEntry()
		newSuperPwdEntry2.SetPlaceHolder("请再次输入")

		// 创建表单对话框
		form := dialog.NewForm("超级密码修改", "确定", "取消", []*widget.FormItem{
			{Text: "原密码", Widget: oldSuperPwdEntry, Required: true},
			{Text: "新密码", Widget: newSuperPwdEntry, Required: true},
			{Text: "确认密码", Widget: newSuperPwdEntry2, Required: true},
		}, func(confirm bool) {
			if confirm {
				// 获取表单数据
				oldSuperPwd := oldSuperPwdEntry.Text
				newSuperPwd := newSuperPwdEntry.Text
				newSuperPwd2 := newSuperPwdEntry2.Text

				// 检查密码
				if oldSuperPwd == "" || newSuperPwd == "" || newSuperPwd2 == "" {
					dialog.ShowError(errors.New("密码不能为空"), w)
					return
				}
				if newSuperPwd != newSuperPwd2 {
					dialog.ShowError(errors.New("两次输入的新密码不一致"), w)
					return
				}
				// 验证旧密码
				var setting global.Setting
				err := db.Instance().Model(&global.Setting{}).First(&setting, "key = ?", "SuperPWD").Error
				if err != nil {
					dialog.ShowError(errors.New("查询超级密码失败"+err.Error()), w)
					return
				}

				tmpSuperPwd, _ := utils.RSADecrypt(privateKey, setting.Value)
				if oldSuperPwd != tmpSuperPwd {
					dialog.ShowError(errors.New("原密码不正确"), w)
					return
				}
				// 更新超级密码
				newSuperPwdEnc, _ := utils.RSAEncrypt(publicKey, newSuperPwd)
				err = db.Instance().Model(&global.Setting{}).Where("key = ?", "SuperPWD").Update("value", newSuperPwdEnc).Error
				if err != nil {
					dialog.ShowError(errors.New("更新超级密码失败"+err.Error()), w)
					return
				}
				log.Instance().Info("修改超级密码成功")
				dialog.ShowInformation("更新超级密码成功", "超级密码已成功更新", w)
			}
		}, w)

		// 设置表单对话框的大小
		form.Resize(fyne.NewSize(300, 200))

		form.Show()
	})
}

func createForgetSuperPwdButton(w fyne.Window) *widget.Button {
	return widget.NewButton("超级密码重置", func() {
		recoverQuestion1 := widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
		recoverAnswer1 := widget.NewEntry()
		recoverAnswer1.SetPlaceHolder("请输入答案")

		recoverQuestion2 := widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
		recoverAnswer2 := widget.NewEntry()
		recoverAnswer2.SetPlaceHolder("请输入答案")

		recoverQuestion3 := widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
		recoverAnswer3 := widget.NewEntry()
		recoverAnswer3.SetPlaceHolder("请输入答案")

		// 创建表单对话框
		form := dialog.NewForm("重置超级密码", "确定", "取消", []*widget.FormItem{
			{Text: "问题1", Widget: recoverQuestion1},
			{Text: "答案1", Widget: recoverAnswer1, Required: true},
			{Text: "问题2", Widget: recoverQuestion2},
			{Text: "答案2", Widget: recoverAnswer2, Required: true},
			{Text: "问题3", Widget: recoverQuestion3},
			{Text: "答案3", Widget: recoverAnswer3, Required: true},
		}, func(confirm bool) {
			if confirm {
				// 获取表单数据
				recoverAnswer1Text := recoverAnswer1.Text
				recoverAnswer2Text := recoverAnswer2.Text
				recoverAnswer3Text := recoverAnswer3.Text

				if recoverAnswer1Text == "" || recoverAnswer2Text == "" || recoverAnswer3Text == "" {
					dialog.ShowError(errors.New("答案不能为空"), w)
					return
				}
				// 验证问题1答案
				var answer1 global.Setting
				err := db.Instance().Model(&global.Setting{}).First(&answer1, "key = ?", "RecoverAnswer1").Error
				if err != nil {
					dialog.ShowError(errors.New("查询问题1答案失败"+err.Error()), w)
					return
				}
				answer1Value, _ := utils.RSADecrypt(privateKey, answer1.Value)
				if recoverAnswer1Text != answer1Value {
					dialog.ShowError(errors.New("问题1答案不正确"), w)
					return
				}
				// 验证问题2答案
				var answer2 global.Setting
				err = db.Instance().Model(&global.Setting{}).First(&answer2, "key = ?", "RecoverAnswer2").Error
				if err != nil {
					dialog.ShowError(errors.New("查询问题2答案失败"+err.Error()), w)
					return
				}
				answer2Value, _ := utils.RSADecrypt(privateKey, answer2.Value)
				if recoverAnswer2Text != answer2Value {
					dialog.ShowError(errors.New("问题2答案不正确"), w)
					return
				}
				// 验证问题3答案
				var answer3 global.Setting
				err = db.Instance().Model(&global.Setting{}).First(&answer3, "key = ?", "RecoverAnswer3").Error
				if err != nil {
					dialog.ShowError(errors.New("查询问题3答案失败"+err.Error()), w)
					return
				}
				answer3Value, _ := utils.RSADecrypt(privateKey, answer3.Value)
				if recoverAnswer3Text != answer3Value {
					dialog.ShowError(errors.New("问题3答案不正确"), w)
					return
				}
				// 重置超级密码
				newSuperPwdEntry := widget.NewPasswordEntry()
				newSuperPwdEntry.SetPlaceHolder("请输入新超级密码")
				newSuperPwdEntry2 := widget.NewPasswordEntry()
				newSuperPwdEntry2.SetPlaceHolder("请再次输入")
				// 创建表单对话框, 二次弹窗
				form := dialog.NewForm("重置超级密码", "确定", "取消", []*widget.FormItem{
					{Text: "新密码", Widget: newSuperPwdEntry, Required: true},
					{Text: "确认密码", Widget: newSuperPwdEntry2, Required: true},
				}, func(confirm bool) {
					if confirm {
						// 获取表单数据
						newSuperPwd := newSuperPwdEntry.Text
						newSuperPwd2 := newSuperPwdEntry2.Text
						if newSuperPwd == "" || newSuperPwd2 == "" {
							dialog.ShowError(errors.New("密码不能为空"), w)
							return
						}
						if newSuperPwd != newSuperPwd2 {
							dialog.ShowError(errors.New("两次输入的新密码不一致"), w)
							return
						}
						// 更新超级密码
						newSuperPwdEnc, _ := utils.RSAEncrypt(publicKey, newSuperPwd)
						err = db.Instance().Model(&global.Setting{}).Where("key = ?", "SuperPWD").Update("value", newSuperPwdEnc).Error
						if err != nil {
							dialog.ShowError(errors.New("更新超级密码失败"+err.Error()), w)
							return
						}
						log.Instance().Info("重置超级密码成功")
						dialog.ShowInformation("更新超级密码成功", "超级密码已成功更新", w)
					}
				}, w)
				// 设置表单对话框的大小
				form.Resize(fyne.NewSize(300, 200))
				// 显示表单对话框
				form.Show()
			}
		}, w)

		// 设置表单对话框的大小
		form.Resize(fyne.NewSize(300, 200))

		// 查询问题, 并渲染
		var question1 global.Setting
		err := db.Instance().Model(&global.Setting{}).First(&question1, "key = ?", "RecoverQuestion1").Error
		if err != nil {
			dialog.ShowError(errors.New("查询问题1失败"+err.Error()), w)
			return
		}

		recoverQuestion1.SetText(question1.Value)
		var question2 global.Setting
		err = db.Instance().Model(&global.Setting{}).First(&question2, "key = ?", "RecoverQuestion2").Error
		if err != nil {
			dialog.ShowError(errors.New("查询问题2失败"+err.Error()), w)
			return
		}
		recoverQuestion2.SetText(question2.Value)
		var question3 global.Setting
		err = db.Instance().Model(&global.Setting{}).First(&question3, "key = ?", "RecoverQuestion3").Error
		if err != nil {
			dialog.ShowError(errors.New("查询问题3失败"+err.Error()), w)
			return
		}
		recoverQuestion3.SetText(question3.Value)

		// 显示表单对话框
		form.Show()
	})
}

func createStatusSuperPwdButton(w fyne.Window) *widget.Button {
	return widget.NewButton("超级密码启停", func() {
		// 创建表单字段
		statusEntry := widget.NewRadioGroup([]string{"是", "否"}, nil)

		// 创建表单对话框
		form := dialog.NewForm("超级密码状态", "确定", "取消", []*widget.FormItem{
			{Text: "启用状态", Widget: statusEntry, Required: true},
		}, func(confirm bool) {
			if confirm {
				// 获取表单数据
				mStatus := statusEntry.Selected
				// 查询原值
				var setting global.Setting
				err := db.Instance().Model(&global.Setting{}).First(&setting, "key = ?", "SuperPwdStatus").Error
				if err != nil {
					dialog.ShowError(errors.New("查询超级密码状态失败"+err.Error()), w)
					return
				}
				// 检查是否修改, 是否需要更新
				if (mStatus == "否" && setting.Value == "false") || (mStatus == "是" && setting.Value == "true") {
					//dialog.ShowInformation("提示", "当前超级密码状态已为: "+mStatus+"\n无需修改", w)
					return
				}
				// 根据选择转换状态值
				var status string
				if mStatus == "否" {
					status = "false"
				} else if mStatus == "是" {
					status = "true"
				}
				// 执行更新操作
				err = db.Instance().Model(&global.Setting{}).Where("key = ?", "SuperPwdStatus").Update("value", status).Error
				if err != nil {
					dialog.ShowError(errors.New("更新超级密码状态失败"+err.Error()), w)
					return
				}
				log.Instance().Info("更新超级密码状态成功")
				dialog.ShowInformation("更新超级密码状态成功", "超级密码启用状态变更为: "+mStatus+"\n请重新启动应用生效", w)
				time.AfterFunc(2000*time.Millisecond, func() {
					fyne.Do(func() {
						w.Close()
					})
				})
			}
		}, w)
		// 设置表单对话框的大小
		form.Resize(fyne.NewSize(300, 200))

		// 查询超级密码状态
		var setting global.Setting
		err := db.Instance().First(&setting, "key = ?", "SuperPwdStatus").Error
		if err != nil {
			dialog.ShowError(errors.New("查询超级密码状态失败"+err.Error()), w)
			return
		}
		if setting.Value == "false" {
			statusEntry.SetSelected("否")
		} else if setting.Value == "true" {
			statusEntry.SetSelected("是")
		}

		form.Show()
	})
}

func createClearDataButton(w fyne.Window) *widget.Button {
	return widget.NewButton("密码数据清空", func() {
		// 清空密码文件
		dialog.ShowConfirm("清空数据", "确定要清空数据吗？", func(confirmed bool) {
			if confirmed {
				// 执行清空操作
				// 清空密码数据
				err := db.Instance().Exec("DELETE FROM " + global.PassWord{}.TableName()).Error
				if err != nil {
					dialog.ShowError(errors.New("清空密码数据失败"+err.Error()), w)
					return
				}
				err = db.Instance().Exec("UPDATE sqlite_sequence SET seq = 0 WHERE name = ?", global.PassWord{}.TableName()).Error
				if err != nil {
					dialog.ShowError(errors.New("重置密码数据计数失败"+err.Error()), w)
					return
				}

				// 清空配置数据
				err = db.Instance().Exec("DELETE FROM " + global.Setting{}.TableName()).Error
				if err != nil {
					dialog.ShowError(errors.New("清空配置数据失败"+err.Error()), w)
					return
				}

				err = db.Instance().Exec("UPDATE sqlite_sequence SET seq = 0 WHERE name = ?", global.Setting{}.TableName()).Error
				if err != nil {
					dialog.ShowError(errors.New("重置配置数据计数失败"+err.Error()), w)
					return
				}
				log.Instance().Info("清空数据成功")
				dialog.ShowInformation("清空数据成功", "密码数据已成功清空，请重新打开应用", w)
				time.AfterFunc(2000*time.Millisecond, func() {
					fyne.Do(func() {
						w.Close()
					})
				})
			}
		}, w)
	})
}

// --------------------------文件操作层--------------------------
func importPwdFromFile(filePath string) (error, int) {
	// 读取密码表并写入密码数据
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		log.Instance().Error("打开文件失败: " + err.Error())
		return err, 0
	}
	// 关闭表
	defer func() { _ = f.Close() }()

	rows, err := f.GetRows("Sheet1")
	if err != nil {
		log.Instance().Error("读取文件失败: " + err.Error())
		return err, 0
	}
	// 遍历行数据
	var count int
	for i, row := range rows {
		if i == 0 {
			continue
		}
		var pwd global.PassWord
		pwd.Title = row[0]
		pwd.Username = row[1]
		pwd.Password = row[2]
		pwd.Category = row[3]
		pwd.Url = row[4]
		pwd.Remark = row[5]

		// 检查表格数据是否为空
		if pwd.Title == "" || pwd.Username == "" || pwd.Password == "" || pwd.Category == "" {
			log.Instance().Error("处理第" + strconv.Itoa(i) + "数据失败: 数据为空")
			continue
		}
		// 密码加密
		pwd.Password, _ = utils.RSAEncrypt(publicKey, pwd.Password)
		// 检查数据是否已存在, 如果存在则更新, 否则插入新数据
		var existPwd global.PassWord
		if db.Instance().Where("title = ? AND username = ? ", pwd.Title, pwd.Username).First(&existPwd).Error == nil {
			// 更新该数据
			err = db.Instance().Model(&global.PassWord{}).Where("id = ?", existPwd.ID).
				Updates(&pwd).Error
			if err != nil {
				log.Instance().Error("更新第" + strconv.Itoa(i) + "数据失败: " + err.Error())
				continue
			}
			count += 1
			log.Instance().Info("更新第" + strconv.Itoa(i) + "数据成功，数据已存在")
			continue
		}
		// 插入新数据
		err = db.Instance().Create(&pwd).Error
		if err != nil {
			log.Instance().Error("插入第" + strconv.Itoa(i) + "数据失败: " + err.Error())
			continue
		}
		count += 1
		log.Instance().Info("插入第" + strconv.Itoa(i) + "数据成功")
	}
	return nil, count
}

func exportPwdToFile(data []global.PassWord, writer fyne.URIWriteCloser) error {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			log.Instance().Error("关闭文件失败: " + err.Error())
			return
		}
	}()
	// 创建一个工作表
	index, err := f.NewSheet("Sheet1")
	if err != nil {
		return err
	}
	// 设置工作簿的默认工作表
	f.SetActiveSheet(index)
	// 绘制表头
	f.SetCellValue("Sheet1", "A1", "标题")
	f.SetCellValue("Sheet1", "B1", "用户名")
	f.SetCellValue("Sheet1", "C1", "密码")
	f.SetCellValue("Sheet1", "D1", "类别")
	f.SetCellValue("Sheet1", "E1", "地址")
	f.SetCellValue("Sheet1", "F1", "备注")
	for i, pwd := range data {
		tmpPwd, _ := utils.RSADecrypt(privateKey, pwd.Password) //解密为明文
		f.SetCellValue("Sheet1", fmt.Sprintf("A%d", i+2), pwd.Title)
		f.SetCellValue("Sheet1", fmt.Sprintf("B%d", i+2), pwd.Username)
		f.SetCellValue("Sheet1", fmt.Sprintf("C%d", i+2), tmpPwd)
		f.SetCellValue("Sheet1", fmt.Sprintf("D%d", i+2), pwd.Category)
		f.SetCellValue("Sheet1", fmt.Sprintf("E%d", i+2), pwd.Url)
		f.SetCellValue("Sheet1", fmt.Sprintf("F%d", i+2), pwd.Remark)
	}

	// 根据指定路径保存文件
	if err = f.SaveAs(writer.URI().Path()); err != nil {
		log.Instance().Error("保存文件失败: " + err.Error())
		return err
	}

	log.Instance().Info("保存文件成功, 文件: " + writer.URI().Path())
	return nil
}
