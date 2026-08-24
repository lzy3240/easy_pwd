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
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"net/url"
	"strings"
	"time"
)

type searchPwdParam struct {
	Title string
	Url   string
	Cate  string
}

var (
	pwdDataList        []global.PassWord
	selectedPwd        global.PassWord
	pwdTable           *widget.Table = nil
	selectedPwdIndex   int           = -1
	questionList       []string      = func() []string { return strings.Split(global.QuestionString, "\n") }()
	publicKey          string        = getSetting("PublicKey")
	privateKey         string        = getSetting("PrivateKey")
	superPwdStatus     bool          = getSetting("SuperPwdStatus") == "true"
	superPwdSetStatus  bool          = getSuperPwdStatus()
	currentCheckStatus bool
)

// -------------------主视图---------------------
func createPasswordView(w fyne.Window, a fyne.App) *fyne.Container {
	// 创建搜索容器
	searchContainer := createPwdSearchContainer(w)

	// 初始化密码列表容器
	pwdTable = createPwdTableContainer()

	// 创建工具容器
	toolContainer := createPwdToolContainer(w, a)

	// top容器
	topContainer := container.NewBorder(
		nil, nil, nil, searchContainer,
	)

	// middle容器
	middleContainer := container.NewBorder(
		nil, nil, nil, nil,
		pwdTable,
	)

	// bottom容器
	bottomContainer := container.NewVBox(
		toolContainer,
	)

	// 初始化全部密码数据
	queryPwdData()

	// 检查系统是否设置超级密码
	if superPwdStatus {
		superPwdSetCheck(w)
	}

	// 主容器
	return container.NewBorder(
		topContainer,
		bottomContainer,
		nil,
		nil,
		middleContainer,
	)
}

// -------------------容器层---------------------
func createPwdSearchContainer(w fyne.Window) *fyne.Container {
	titleEntry := widget.NewEntry()
	titleEntry.SetPlaceHolder("请输入标题")
	searchContainer1 := container.NewHBox(
		widget.NewLabel("标题"),
		container.NewGridWrap(fyne.NewSize(120, 40), titleEntry),
	)

	urlEntry := widget.NewEntry()
	urlEntry.SetPlaceHolder("请输入地址")
	searchContainer2 := container.NewHBox(
		widget.NewLabel("地址"),
		container.NewGridWrap(fyne.NewSize(120, 40), urlEntry),
	)

	cateEntry := widget.NewSelect([]string{"常规", "服务器", "数据库", "网站", "邮箱", "网银"}, nil)
	cateEntry.PlaceHolder = "请选择" //设置占位符
	searchContainer3 := container.NewHBox(
		widget.NewLabel("类型"),
		container.NewGridWrap(fyne.NewSize(120, 40), cateEntry),
	)

	// 回车搜索
	titleEntry.OnSubmitted = func(text string) {
		searchParam := &searchPwdParam{
			Title: titleEntry.Text,
			Url:   urlEntry.Text,
			Cate:  cateEntry.Selected,
		}
		searchPwdData(w, searchParam, pwdTable)
	}
	urlEntry.OnSubmitted = func(text string) {
		searchParam := &searchPwdParam{
			Title: titleEntry.Text,
			Url:   urlEntry.Text,
			Cate:  cateEntry.Selected,
		}
		searchPwdData(w, searchParam, pwdTable)
	}
	cateEntry.OnChanged = func(text string) {
		searchParam := &searchPwdParam{
			Title: titleEntry.Text,
			Url:   urlEntry.Text,
			Cate:  cateEntry.Selected,
		}
		searchPwdData(w, searchParam, pwdTable)
	}

	// 搜索按钮
	searchButton := widget.NewButton("搜索", func() {
		searchParam := &searchPwdParam{
			Title: titleEntry.Text,
			Url:   urlEntry.Text,
			Cate:  cateEntry.Selected,
		}
		searchPwdData(w, searchParam, pwdTable)
	})
	searchButton.Importance = widget.HighImportance
	searchButton.SetIcon(theme.SearchIcon())

	resetButton := widget.NewButton("重置", func() {
		titleEntry.SetText("")
		urlEntry.SetText("")
		cateEntry.ClearSelected()
		searchPwdData(w, &searchPwdParam{}, pwdTable)
	})
	resetButton.Importance = widget.HighImportance
	resetButton.SetIcon(theme.SearchReplaceIcon())

	grid := container.NewGridWithColumns(2,
		searchButton,
		resetButton,
	)

	searchContainer := container.NewHBox(
		searchContainer1,
		searchContainer2,
		searchContainer3,
		grid,
	)
	return searchContainer
}

func createPwdTableContainer() *widget.Table {
	// 创建一个带有表头的表格
	table := widget.NewTable(
		func() (int, int) { return len(pwdDataList), 7 }, // 修正行列数，n行7列
		func() fyne.CanvasObject {
			// 创建一个标签用于显示密码记录
			lab := widget.NewLabel("")

			// 设置标签超出部分的处理方式
			lab.Truncation = fyne.TextTruncateEllipsis

			return lab
		},
		// 渲染全部数据到表格
		setPwdDataToTable,
	)

	// 启用表头
	table.ShowHeaderRow = true

	// 自定义表头内容的创建
	table.CreateHeader = func() fyne.CanvasObject {
		return widget.NewLabel("Header Template")
	}

	// 更新自定义表头内容
	table.UpdateHeader = func(id widget.TableCellID, template fyne.CanvasObject) {
		// 根据列索引设置表头标签的文本内容
		switch id.Col {
		case 0:
			template.(*widget.Label).SetText("标题")
		case 1:
			template.(*widget.Label).SetText("用户名")
		case 2:
			template.(*widget.Label).SetText("密码")
		case 3:
			template.(*widget.Label).SetText("类型")
		case 4:
			template.(*widget.Label).SetText("地址")
		case 5:
			template.(*widget.Label).SetText("修改时间")
		case 6:
			template.(*widget.Label).SetText("备注")
		}

		// 设置表头标签的样式
		template.(*widget.Label).TextStyle = fyne.TextStyle{Bold: true}

		// 设置表头标签的对齐方式 - 左对齐
		template.(*widget.Label).Alignment = fyne.TextAlignLeading
	}

	// 设置标题列的宽度
	table.SetColumnWidth(0, 120)
	table.SetColumnWidth(1, 100)
	table.SetColumnWidth(2, 150)
	table.SetColumnWidth(3, 70)
	table.SetColumnWidth(4, 240)
	table.SetColumnWidth(5, 150)
	table.SetColumnWidth(6, 90)

	// 设置表格的选中事件
	table.OnSelected = func(id widget.TableCellID) {
		// 检查是否点击了表头行
		if id.Row == -1 {
			return
		}

		// 检查是否在有效的行上点击
		if id.Row >= len(pwdDataList) {
			return
		}

		// 获取选中的任务
		selectedPwd = pwdDataList[id.Row]
		selectedPwdIndex = id.Row
		// 刷新表格
		table.Refresh()
	}

	// 设置表格取消选中事件
	table.OnUnselected = func(id widget.TableCellID) {
		// 检查是否在表头行上点击
		if id.Row == -1 {
			return
		}

		// 检查是否在有效的行上点击
		if id.Row >= len(pwdDataList) {
			return
		}

		// 清空选中的任务
		selectedPwd = global.PassWord{}
		selectedPwdIndex = -1
		// 刷新表格
		table.Refresh()
	}

	return table
}

func createPwdToolContainer(w fyne.Window, a fyne.App) *fyne.Container {
	return container.NewGridWithColumns(8,
		addPwdButton(pwdTable, w),
		editPwdButton(pwdTable, w),
		deletePwdButton(pwdTable, w),
		refreshPwdButton(pwdTable, w),
		decryptPwdButton(pwdTable, w),
		openUrlButton(pwdTable, w, a),
		copyUserNameButton(pwdTable, w, a),
		copyPwdButton(pwdTable, w, a),
	)
}

// -------------------按钮层---------------------
func addPwdButton(table *widget.Table, w fyne.Window) fyne.CanvasObject {
	// 创建表单字段
	titleEntry := widget.NewEntry()
	titleEntry.SetPlaceHolder("请输入标题")

	userNameEntry := widget.NewEntry()
	userNameEntry.SetPlaceHolder("请输入用户名")

	pwdEntry := widget.NewPasswordEntry()
	pwdEntry.SetPlaceHolder("请输入密码")

	cateEntry := widget.NewSelect([]string{"常规", "服务器", "数据库", "网站", "邮箱", "网银"}, nil)
	cateEntry.PlaceHolder = "请选择类型"

	urlEntry := widget.NewEntry()
	urlEntry.SetPlaceHolder("请输入地址")

	remarkEntry := widget.NewEntry()
	remarkEntry.SetPlaceHolder("请输入备注")

	// 创建表单对话框
	form := dialog.NewForm("添加密码", "确定", "取消", []*widget.FormItem{
		{Text: "标题", Widget: titleEntry, Required: true},
		{Text: "用户名", Widget: userNameEntry, Required: true},
		{Text: "密码", Widget: pwdEntry, Required: true},
		{Text: "类型", Widget: cateEntry, Required: true},
		{Text: "地址", Widget: urlEntry},
		{Text: "备注", Widget: remarkEntry},
	}, func(confirm bool) {
		// 如果用户点击确认则执行
		if confirm {
			// 获取表单字段的值
			title := titleEntry.Text
			userName := userNameEntry.Text
			pwd := pwdEntry.Text
			cate := cateEntry.Selected
			url := urlEntry.Text
			remark := remarkEntry.Text

			if title == "" || userName == "" || pwd == "" || cate == "" {
				dialog.ShowError(errors.New("请填写完整的表单"), w)
				return
			}

			// 对密码进行公钥加密
			tmpPwd, err := utils.RSAEncrypt(publicKey, pwd)
			if err != nil {
				dialog.ShowError(errors.New("密码加密失败"), w)
				return
			}
			//log.Instance().Info("密码加密成功", "密文", tmpPwd)
			// 创建一个密码记录
			pwdData := global.PassWord{
				Title:    title,
				Username: userName,
				Password: tmpPwd, //加密
				Category: cate,
				Url:      url,
				Remark:   remark,
			}

			// 唯一性验证
			var existPwd global.PassWord
			err = db.Instance().Where("title = ? AND username = ?", title, userName).First(&existPwd).Error
			if err == nil && existPwd.ID != 0 {
				dialog.ShowError(errors.New("该标题和用户的密码已存在"), w)
				return
			}

			// 将密码记录保存到数据库
			err = db.Instance().Create(&pwdData).Error
			if err != nil {
				log.Instance().Error("添加密码失败: " + err.Error())
				dialog.ShowError(errors.New("添加密码失败"), w)
				return
			}

			// 刷新密码信息表格
			queryPwdData()

			// 显示成功对话框
			fyne.Do(func() {
				log.Instance().Info("添加密码成功: " + pwdData.Title + ":" + pwdData.Username)
				dialog.ShowInformation("提示", "添加密码成功", w)
			})

			// 清空表单
			//titleEntry.SetText("")
			//userNameEntry.SetText("")
			//pwdEntry.SetText("")
			//cateEntry.ClearSelected()
			//urlEntry.SetText("")
			//remarkEntry.SetText("")

			// 取消选中的行
			table.UnselectAll()
		}

		// 保存后清空表单
		titleEntry.SetText("")
		userNameEntry.SetText("")
		pwdEntry.SetText("")
		cateEntry.ClearSelected()
		urlEntry.SetText("")
		remarkEntry.SetText("")
	}, w)

	// 设置表单对话框的大小
	form.Resize(fyne.NewSize(600, 400))

	// 创建一个添加按钮
	addButton := widget.NewButton("添加", func() {
		form.Show()
	})

	// 设置添加按钮的重要性
	addButton.Importance = widget.HighImportance

	// 设置添加按钮的图标,使用内置的添加图标
	addButton.SetIcon(theme.ContentAddIcon())

	return addButton
}

func editPwdButton(table *widget.Table, w fyne.Window) fyne.CanvasObject {
	// 创建表单字段
	titleEntry := widget.NewEntry()
	userNameEntry := widget.NewEntry()
	pwdEntry := widget.NewPasswordEntry()
	cateEntry := widget.NewSelect([]string{"常规", "服务器", "数据库", "网站", "邮箱", "网银"}, nil)
	urlEntry := widget.NewEntry()
	remarkEntry := widget.NewEntry()

	// 创建表单对话框
	form := dialog.NewForm("编辑密码", "确定", "取消", []*widget.FormItem{
		{Text: "标题", Widget: titleEntry, Required: true},
		{Text: "用户名", Widget: userNameEntry, Required: true},
		{Text: "密码", Widget: pwdEntry, Required: true},
		{Text: "类型", Widget: cateEntry, Required: true},
		{Text: "地址", Widget: urlEntry},
		{Text: "备注", Widget: remarkEntry},
	}, func(confirm bool) {
		// 如果用户点击确认则执行
		if confirm {
			// 获取表单字段的值
			title := titleEntry.Text
			userName := userNameEntry.Text
			pwd := pwdEntry.Text
			cate := cateEntry.Selected
			url := urlEntry.Text
			remark := remarkEntry.Text

			if title == "" || userName == "" || pwd == "" || cate == "" {
				dialog.ShowError(errors.New("请填写完整的表单"), w)
				return
			}

			//// 检查密码是否修改
			//var tmpPwd string
			//var err error
			//if !utils.IsAllAsterisk(pwd) {
			//	// 对密码进行公钥加密
			//	tmpPwd, err = utils.RSAEncrypt(publicKey, pwd)
			//	if err != nil {
			//		dialog.ShowError(errors.New("密码加密失败"), w)
			//		return
			//	}
			//} else {
			//	tmpPwd = "" // 未修改, 利用gorm不更新零值
			//}
			tmpPwd, err := utils.RSAEncrypt(publicKey, pwd)
			if err != nil {
				dialog.ShowError(errors.New("密码加密失败"), w)
				return
			}
			// 创建一个密码记录
			pwdData := global.PassWord{
				Title:    title,
				Username: userName,
				Password: tmpPwd, //已加密或零值
				Category: cate,
				Url:      url,
				Remark:   remark,
			}
			// 唯一性验证
			var existPwd global.PassWord
			err = db.Instance().Where("title = ? AND username = ? AND id != ?", title, userName, selectedPwd.ID).First(&existPwd).Error
			if err == nil && existPwd.ID != 0 {
				dialog.ShowError(errors.New("该标题和用户的密码已存在"), w)
				return
			}

			// 更新密码
			err = db.Instance().Model(&global.PassWord{}).
				Where("id = ?", selectedPwd.ID).
				Updates(&pwdData).Error
			if err != nil {
				log.Instance().Error("编辑密码失败: " + err.Error())
				dialog.ShowError(errors.New("编辑密码失败"), w)
				return
			}

			// 刷新密码信息表格
			queryPwdData()

			// 显示成功对话框
			fyne.Do(func() {
				log.Instance().Info("编辑密码成功: " + selectedPwd.Title + ":" + selectedPwd.Username)
				dialog.ShowInformation("提示", "编辑密码成功", w)
			})

			// 清空表单
			//titleEntry.SetText("")
			//userNameEntry.SetText("")
			//pwdEntry.SetText("")
			//cateEntry.ClearSelected()
			//urlEntry.SetText("")
			//remarkEntry.SetText("")

			// 取消选中的行
			table.UnselectAll()
		}
		// 保存后清空表单
		titleEntry.SetText("")
		userNameEntry.SetText("")
		pwdEntry.SetText("")
		cateEntry.ClearSelected()
		urlEntry.SetText("")
		remarkEntry.SetText("")
	}, w)

	// 设置表单对话框的大小
	form.Resize(fyne.NewSize(600, 400))

	// 创建一个编辑按钮
	editButton := widget.NewButton("编辑", func() {
		// 检查是否有选中的行
		if selectedPwd.ID == 0 {
			// 显示错误对话框
			fyne.Do(func() {
				dialog.ShowError(fmt.Errorf("请先选择要编辑的密码"), w)
			})
			return
		}

		// 解密密码
		if utils.IsBase64(selectedPwd.Password) {
			selectedPwd.Password, _ = utils.RSADecrypt(privateKey, selectedPwd.Password)
		}

		// 设置表单字段的值
		titleEntry.SetText(selectedPwd.Title)
		userNameEntry.SetText(selectedPwd.Username)
		pwdEntry.SetText(selectedPwd.Password)
		cateEntry.SetSelected(selectedPwd.Category)
		urlEntry.SetText(selectedPwd.Url)
		remarkEntry.SetText(selectedPwd.Remark)

		// 打开编辑对话框
		form.Show()
	})

	// 设置按钮的重要性
	editButton.Importance = widget.HighImportance

	// 设置编辑按钮的图标
	editButton.SetIcon(theme.DocumentCreateIcon())

	return editButton
}

func deletePwdButton(table *widget.Table, w fyne.Window) fyne.CanvasObject {
	// 创建一个删除按钮
	deleteButton := widget.NewButton("删除", func() {
		// 检查是否有选中的行
		if selectedPwd.ID == 0 {
			// 显示错误对话框
			fyne.Do(func() {
				dialog.ShowError(fmt.Errorf("请先选择要删除的密码"), w)
			})
			return
		} else {
			fyne.Do(func() {
				dialog.ShowConfirm("删除密码", "确定要删除密码 "+selectedPwd.Title+" 吗？", func(confirmed bool) {
					if confirmed {
						// 删除密码的逻辑
						err := db.Instance().Model(&global.PassWord{}).
							Where("id = ?", selectedPwd.ID).
							Delete(&global.PassWord{}).
							Error
						if err != nil {
							log.Instance().Error("删除密码失败: " + err.Error())
							return
						}

						log.Instance().Info("删除密码成功: " + selectedPwd.Title)

						// 刷新表格数据
						queryPwdData()

						// 显示成功对话框
						fyne.Do(func() {
							log.Instance().Info("删除密码成功: " + selectedPwd.Title + ":" + selectedPwd.Username)
							dialog.ShowInformation("删除密码成功", fmt.Sprintf("已成功删除密码 %s", selectedPwd.Title), w)
						})

						// 清空当前选中的行
						selectedPwd = global.PassWord{}

						// 取消选择所有行
						table.UnselectAll()
					}
				}, w)
			})
		}
	})

	// 设置删除按钮的重要性
	deleteButton.Importance = widget.DangerImportance

	// 设置删除按钮的图标,使用内置的删除图标
	deleteButton.SetIcon(theme.DeleteIcon())

	return deleteButton
}

func decryptPwdButton(table *widget.Table, w fyne.Window) fyne.CanvasObject {
	//selectedPwd.ID
	decryptButton := widget.NewButton("解密", func() {
		// 解密密码的逻辑
		// 检查是否有选中的行
		if selectedPwd.ID == 0 {
			// 显示错误对话框
			fyne.Do(func() {
				dialog.ShowError(fmt.Errorf("请先选择要解密的密码"), w)
			})
			return
		}

		// 修改密码数据列表中选中的密码数据的密码字段
		plaintext, err := utils.RSADecrypt(privateKey, pwdDataList[selectedPwdIndex].Password)
		if err != nil {
			// 显示错误对话框
			fyne.Do(func() {
				dialog.ShowError(fmt.Errorf("密码解密失败或密文错误"), w)
			})
			return
		} else {
			pwdDataList[selectedPwdIndex].Password = plaintext
		}

		// 刷新表格数据
		flushPwdDataToTable(table, pwdDataList)
	})

	decryptButton.Importance = widget.WarningImportance
	decryptButton.SetIcon(theme.VisibilityIcon())

	return decryptButton
}

func refreshPwdButton(table *widget.Table, w fyne.Window) fyne.CanvasObject {
	refreshButton := widget.NewButton("刷新", func() {
		// 取消选择所有行
		table.UnselectAll()

		// 清空当前选中的行
		selectedPwd = global.PassWord{}

		// 查询最新的密码记录
		queryPwdData()

		// 检查查询结果
		if len(pwdDataList) == 0 {
			// 显示提示对话框
			dialog.ShowInformation("提示", "没有找到任何密码数据", w)
			return
		}
	})
	// 设置刷新按钮的重要性
	refreshButton.Importance = widget.HighImportance

	// 设置刷新按钮的图标,使用内置的刷新图标
	refreshButton.SetIcon(theme.ViewRefreshIcon())

	return refreshButton
}

func copyUserNameButton(table *widget.Table, w fyne.Window, a fyne.App) fyne.CanvasObject {
	copyButton := widget.NewButton("复制用户", func() {
		if selectedPwd.ID == 0 {
			dialog.ShowError(errors.New("请先选择要复制的密码"), w)
			return
		}
		// 获取剪贴板对象
		clipboard := a.Clipboard()

		// 清空剪切板
		clipboard.SetContent("")

		// 复制密码到剪切板
		clipboard.SetContent(selectedPwd.Username)

		// 显示成功对话框
		fyne.Do(func() {
			dialog.ShowInformation("复制成功", fmt.Sprintf("已复制用户名 %s ", selectedPwd.Username), w)
		})

		// 取消选择所有行
		table.UnselectAll()
	})

	copyButton.SetIcon(theme.ContentCutIcon())
	copyButton.Importance = widget.HighImportance
	return copyButton
}

func copyPwdButton(table *widget.Table, w fyne.Window, a fyne.App) fyne.CanvasObject {
	copyButton := widget.NewButton("复制密码", func() {
		if selectedPwd.ID == 0 {
			dialog.ShowError(errors.New("请先选择要复制的密码"), w)
			return
		}

		pwd, err := utils.RSADecrypt(privateKey, selectedPwd.Password)
		if err != nil {
			dialog.ShowError(errors.New("密码解密失败或密文错误"), w)
			return
		}

		// 获取剪贴板对象
		clipboard := a.Clipboard()

		// 清空剪切板
		clipboard.SetContent("")

		// 复制密码到剪切板
		clipboard.SetContent(pwd)

		// 显示成功对话框
		fyne.Do(func() {
			dialog.ShowInformation("复制成功", fmt.Sprintf("已复制用户名 %s 的密码", selectedPwd.Username), w)
		})

		// 取消选择所有行
		table.UnselectAll()

	})

	copyButton.SetIcon(theme.ContentCopyIcon())
	copyButton.Importance = widget.HighImportance
	return copyButton
}

func openUrlButton(table *widget.Table, w fyne.Window, a fyne.App) fyne.CanvasObject {
	openButton := widget.NewButtonWithIcon("打开地址", theme.MailSendIcon(), func() {
		if selectedPwd.ID == 0 {
			dialog.ShowError(errors.New("请先选择要打开的密码"), w)
			return
		}
		// 判断是否是有效的URL
		if !utils.IsUrlRegex(selectedPwd.Url) {
			dialog.ShowError(errors.New("无效的URL: "+selectedPwd.Url), w)
			return
		}
		parsedURL, err := url.Parse(selectedPwd.Url)
		if err != nil {
			dialog.ShowError(errors.New("无效的URL: "+selectedPwd.Url), w)
			return
		}
		if err = a.OpenURL(parsedURL); err != nil {
			dialog.ShowError(errors.New("打开URL失败: "+selectedPwd.Url), w)
			return
		}
	})
	//openButton.SetIcon(theme.MailSendIcon())
	openButton.Importance = widget.HighImportance
	return openButton
}

// -------------------业务层---------------------
func searchPwdData(w fyne.Window, entry *searchPwdParam, table *widget.Table) {
	// 如果输入框为空，则默认显示所有密码数据
	if entry.Title == "" && entry.Cate == "" && entry.Url == "" {
		queryPwdData()
		// 检查查询结果
		if len(pwdDataList) == 0 {
			dialog.ShowInformation("提示", "没有找到任何记录", w)
			return
		}
		// 取消选择所有行
		table.UnselectAll()
		return
	}

	// 执行查询, 获取匹配的密码数据
	var pwdList []global.PassWord
	DB := db.Instance().Model(&global.PassWord{})
	if entry.Title != "" {
		DB = DB.Where("title like ?", fmt.Sprintf("%%%s%%", entry.Title))
	}
	if entry.Cate != "" {
		DB = DB.Where("category like ?", fmt.Sprintf("%%%s%%", entry.Cate))
	}
	if entry.Url != "" {
		DB = DB.Where("url like ?", fmt.Sprintf("%%%s%%", entry.Url))
	}

	err := DB.Find(&pwdList).Error
	if err != nil {
		log.Instance().Error("查询密码信息数据失败: " + err.Error())
		dialog.ShowError(errors.New("查询密码信息数据失败"), w)
		return
	}
	// 取消选择所有行
	table.UnselectAll()

	// 刷新表格数据时使用局部变量
	flushPwdDataToTable(table, pwdList)
	// 显示对话框
	if len(pwdList) == 0 {
		dialog.ShowInformation("提示", "没有找到匹配的记录", w)
	}
}

func setPwdDataToTable(id widget.TableCellID, obj fyne.CanvasObject) {
	// 检查行索引是否在有效范围内
	if id.Row < 0 || id.Row >= len(pwdDataList) || pwdDataList == nil {
		obj.(*widget.Label).SetText("") // 如果行索引无效或切片未初始化，设置为空字符串
		return
	}

	// 如果切片为空，则设置为空字符串
	if len(pwdDataList) == 0 {
		obj.(*widget.Label).SetText("")
		return
	}

	// 获取单元格的标签
	label := obj.(*widget.Label)

	// 根据行和列设置标签的文本内容
	if id.Row >= 0 && id.Row < len(pwdDataList) && pwdDataList != nil {
		// 根据行索引获取对应的密码
		data := pwdDataList[id.Row]

		if id.Row == selectedPwdIndex {
			label.Importance = widget.DangerImportance // 高亮显示选中的行
			// 根据列索引设置标签的文本内容
			renderPwdDataToText(id, label, data)
		} else {
			label.Importance = widget.MediumImportance // 正常显示其他行
			// 根据列索引设置标签的文本内容
			renderPwdDataToText(id, label, data)
		}
	} else {
		label.SetText("") // 如果超出行范围，设置为空字符串
	}
}

func renderPwdDataToText(id widget.TableCellID, label *widget.Label, data global.PassWord) {
	switch id.Col {
	case 0: // 标题
		label.SetText(data.Title)
	case 1: // 用户名
		label.SetText(data.Username)
	case 2: // 密码
		label.SetText(data.Password)
	case 3: // 类别
		label.SetText(data.Category)
	case 4: // 地址
		label.SetText(data.Url)
	case 5: // 更新时间
		label.SetText(data.UpdatedAt.Format("2006-01-02 15:04:05")) // Fixed: Changed data.UpdatedAt to string
	case 6: // 备注
		label.SetText(data.Remark)
	default:
		label.SetText("") // 如果超出列范围，设置为空字符串
	}
}

func queryPwdData() {
	// 查询密码信息数据
	var pwdList []global.PassWord
	err := db.Instance().Model(&global.PassWord{}).Find(&pwdList).Error
	if err != nil {
		log.Instance().Error("查询密码数据失败: " + err.Error())
	}

	// 刷新密码数据表格
	flushPwdDataToTable(pwdTable, pwdList)
}

func flushPwdDataToTable(table *widget.Table, data []global.PassWord) {
	pwdDataList = data
	// 刷新表格
	if table != nil {
		// 如果表格未初始化, 则不刷新
		fyne.Do(func() {
			table.Refresh()
		})
	}
}

// -------------------辅助工具---------------------
func getSetting(key string) string {
	var setting global.Setting
	err := db.Instance().Model(&global.Setting{}).Where("key = ?", key).First(&setting).Error
	if err != nil {
		return ""
	}
	return setting.Value
}

func getSuperPwdStatus() bool {
	// 检查超级密码是否存在
	var setting global.Setting
	err := db.Instance().Model(&global.Setting{}).Where("key = ?", "SuperPWD").First(&setting).Error
	if err != nil {
		log.Instance().Error("查询超级密码失败: " + err.Error())
	}
	return setting.ID != 0 && setting.Value != ""
}

// -------------------首次检查---------------------
// 超级密码设置状态检查
func superPwdSetCheck(w fyne.Window) {
	// 检查超级密码是否存在
	if superPwdSetStatus == false {
		// 不存在则弹出创建超级密码的对话框
		// 创建表单字段
		superPwdEntry := widget.NewPasswordEntry()
		superPwdEntry.SetPlaceHolder("请输入超级密码")

		repeatPwdEntry := widget.NewPasswordEntry()
		repeatPwdEntry.SetPlaceHolder("请再次输入超级密码")

		recoverQuestion1Entry := widget.NewSelect(questionList, nil)
		recoverQuestion1Entry.PlaceHolder = "请选择第一个恢复问题"
		recoverAnswer1Entry := widget.NewEntry()
		recoverAnswer1Entry.SetPlaceHolder("请输入答案")

		recoverQuestion2Entry := widget.NewSelect(questionList, nil)
		recoverQuestion2Entry.PlaceHolder = "请选择第二个恢复问题"
		recoverAnswer2Entry := widget.NewEntry()
		recoverAnswer2Entry.SetPlaceHolder("请输入答案")

		recoverQuestion3Entry := widget.NewSelect(questionList, nil)
		recoverQuestion3Entry.PlaceHolder = "请选择第三个恢复问题"
		recoverAnswer3Entry := widget.NewEntry()
		recoverAnswer3Entry.SetPlaceHolder("请输入答案")

		// 创建对话
		form := dialog.NewForm("添加密码", "确定", "取消", []*widget.FormItem{
			{Text: "超级密码", Widget: superPwdEntry, Required: true},
			{Text: "重复超级密码", Widget: repeatPwdEntry, Required: true},
			{Text: "找回问题一", Widget: recoverQuestion1Entry, Required: true},
			{Text: "答案", Widget: recoverAnswer1Entry, Required: true},
			{Text: "找回问题二", Widget: recoverQuestion2Entry, Required: true},
			{Text: "答案", Widget: recoverAnswer2Entry, Required: true},
			{Text: "找回问题三", Widget: recoverQuestion3Entry, Required: true},
			{Text: "找答案", Widget: recoverAnswer3Entry, Required: true},
		}, func(confirm bool) {
			if confirm {
				superPwd := superPwdEntry.Text
				repeatPwd := repeatPwdEntry.Text
				recoverQuestion1 := recoverQuestion1Entry.Selected
				recoverAnswer1 := recoverAnswer1Entry.Text
				recoverQuestion2 := recoverQuestion2Entry.Selected
				recoverAnswer2 := recoverAnswer2Entry.Text
				recoverQuestion3 := recoverQuestion3Entry.Selected
				recoverAnswer3 := recoverAnswer3Entry.Text

				if superPwd == "" || repeatPwd == "" || recoverQuestion1 == "" || recoverAnswer1 == "" || recoverQuestion2 == "" || recoverAnswer2 == "" || recoverQuestion3 == "" || recoverAnswer3 == "" {
					dialog.ShowError(errors.New("请填写完整的超级密码信息"), w)
					// 延迟关闭对话并重载超级密码设置对话框
					time.AfterFunc(2000*time.Millisecond, func() {
						fyne.Do(func() {
							superPwdSetCheck(w)
						})
					})
				} else {
					if superPwd != repeatPwd {
						dialog.ShowError(errors.New("两次输入的超级密码不一致"), w)
						// 延迟关闭对话并重载超级密码设置对话框
						time.AfterFunc(2000*time.Millisecond, func() {
							fyne.Do(func() {
								superPwdSetCheck(w)
							})
						})
					} else {
						// 加密, 使用公钥加密
						superPwd, _ = utils.RSAEncrypt(publicKey, superPwd)
						recoverAnswer1, _ = utils.RSAEncrypt(publicKey, recoverAnswer1)
						recoverAnswer2, _ = utils.RSAEncrypt(publicKey, recoverAnswer2)
						recoverAnswer3, _ = utils.RSAEncrypt(publicKey, recoverAnswer3)

						// 保存超级密码和恢复问题到数据库
						// 删除已有的key, 防止因空值重复创建key
						db.Instance().Where("key = ?", "SuperPWD").Delete(&global.Setting{})
						db.Instance().Where("key = ?", "RecoverQuestion1").Delete(&global.Setting{})
						db.Instance().Where("key = ?", "RecoverAnswer1").Delete(&global.Setting{})
						db.Instance().Where("key = ?", "RecoverQuestion2").Delete(&global.Setting{})
						db.Instance().Where("key = ?", "RecoverAnswer2").Delete(&global.Setting{})
						db.Instance().Where("key = ?", "RecoverQuestion3").Delete(&global.Setting{})
						db.Instance().Where("key = ?", "RecoverAnswer3").Delete(&global.Setting{})
						// 创建新的key
						db.Instance().Create(&global.Setting{Key: "SuperPWD", Value: superPwd})
						db.Instance().Create(&global.Setting{Key: "RecoverQuestion1", Value: recoverQuestion1})
						db.Instance().Create(&global.Setting{Key: "RecoverAnswer1", Value: recoverAnswer1})
						db.Instance().Create(&global.Setting{Key: "RecoverQuestion2", Value: recoverQuestion2})
						db.Instance().Create(&global.Setting{Key: "RecoverAnswer2", Value: recoverAnswer2})
						db.Instance().Create(&global.Setting{Key: "RecoverQuestion3", Value: recoverQuestion3})
						db.Instance().Create(&global.Setting{Key: "RecoverAnswer3", Value: recoverAnswer3})

						// 显示成功对话框
						fyne.Do(func() {
							dialog.ShowInformation("提示", "超级密码设置成功", w)
							// 延迟关闭对话并打开验证状态检查对话框
							time.AfterFunc(2000*time.Millisecond, func() {
								fyne.Do(func() {
									superPwdVerifyCheck(w)
								})
							})
						})
					}
				}
			} else {
				dialog.ShowError(errors.New("第一次使用请进行超级密码设置"), w)
				// 延迟关闭对话并重载超级密码设置对话框
				time.AfterFunc(2000*time.Millisecond, func() {
					fyne.Do(func() {
						superPwdSetCheck(w)
					})
				})
			}
		}, w)
		// 设置表单对话框的大小
		form.Resize(fyne.NewSize(600, 400))
		// 显示表单对话框
		form.Show()
	} else {
		// 已经设置过超级密码，则直接进行验证状态检查
		superPwdVerifyCheck(w)
	}
}

// 验证超级密码状态检查
func superPwdVerifyCheck(w fyne.Window) {
	// 检查当前是否验证过超级密码
	if currentCheckStatus == false {
		// 创建表单字段
		superPwdEntry := widget.NewPasswordEntry()
		superPwdEntry.SetPlaceHolder("请输入超级密码")
		// 创建对话
		form := dialog.NewForm("验证超级密码", "确定", "取消", []*widget.FormItem{
			{Text: "超级密码", Widget: superPwdEntry, Required: true},
		}, func(confirm bool) {
			if confirm {
				superPwd := superPwdEntry.Text
				if superPwd == "" {
					dialog.ShowError(errors.New("请输入超级密码"), w)
					// 延迟关闭对话并重载打开状态检查对话框
					time.AfterFunc(2000*time.Millisecond, func() {
						fyne.Do(func() {
							superPwdVerifyCheck(w)
						})
					})
				} else {
					superPwdDB, _ := utils.RSADecrypt(privateKey, getSetting("SuperPWD"))
					if superPwd == superPwdDB {
						// 显示成功对话框
						fyne.Do(func() {
							dialog.ShowInformation("提示", "超级密码验证成功", w)
						})
						currentCheckStatus = true
					} else {
						dialog.ShowError(errors.New("超级密码验证失败"), w)
						// 延迟关闭对话并重载打开状态检查对话框
						time.AfterFunc(2000*time.Millisecond, func() {
							fyne.Do(func() {
								superPwdVerifyCheck(w)
							})
						})
					}
				}
			} else {
				// 点击取消时
				dialog.ShowError(errors.New("请验证超级密码"), w)
				// 延迟关闭对话并重载打开状态检查对话框
				time.AfterFunc(2000*time.Millisecond, func() {
					fyne.Do(func() {
						superPwdVerifyCheck(w)
					})
				})
			}
		}, w)
		// 设置表单对话框的大小
		form.Resize(fyne.NewSize(300, 200))
		// 显示表单对话框
		form.Show()
	}
}
