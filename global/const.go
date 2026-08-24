package global

import _ "embed"

const (
	DefProjectName   = "easy_pwd"
	DefProjectNameCN = "简密"
	DefRuntimeDir    = "runtime"
	LogFileDir       = "logs"
	DBFileName       = "easy_pwd"
	Version          = "v1.0.0"
	ProjectURL       = "https://github.com/lzy3240/easy_pwd"
	Author           = "zhenyi.law"
	Contact          = "lzy3240@qq.com"
	QuestionString   = "你最喜欢的颜色是什么？\n你第一只宠物的名字是什么？\n你家乡的名字是什么？\n你最喜欢的一部电影是什么？\n你最喜欢的一首歌是什么？\n你最喜欢的一位老师的名字是什么？\n你第一所学校的名称是什么？"
)

//go:embed icon.png
var Icon []byte //包含图标
